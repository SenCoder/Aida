#include <sys/types.h> //数据类型定义
#include <sys/socket.h>//提供socket函数及数据结构
#include <netinet/in.h>//定义数据结构sockaddr_in
#include <arpa/inet.h>//提供IP地址转换函数
#include <netdb.h>//提供设置及获取域名的函数
#include <sys/ioctl.h>//提供对I/O控制的函数
#include <sys/poll.h>//提供socket等待测试机制的函数
#include <errno.h>
#include <fcntl.h>
#include <assert.h>
#include <sys/un.h>
#include <stddef.h>
#include <unistd.h> //提供close
#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>

#define SOCKET_PORT (1024)
#define SOCKET_IS_BLOCKING (1)  //定义是否阻塞发送

// 16 * 16  
static int base64_decode_map[256] = {  
    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, // 0   - 15  
    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, // 16  - 31  
    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 62, -1, -1, -1, 63, // 32  - 47  
    52, 53, 54, 55, 56, 57, 58, 59, 60, 61, -1, -1, -1, -1, -1, -1, // 48  - 63  
    -1,  0,  1,  2,  3,  4,  5,  6,  7,  8,  9, 10, 11, 12, 13, 14, // 64  - 79  
    15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, -1, -1, -1, -1, -1, // 80  - 95  
    -1, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, // 96  - 111  
    41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, -1, -1, -1, -1, -1, // 112 - 127  
    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, // 128 - 143  
    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, // 144 - 159   
    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, // 160 - 175  
    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, // 176 - 191  
    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, // 192 - 207  
    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, // 208 - 223  
    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, // 224 - 239  
    -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, // 240 - 255  
}; 

static char base64_index[65] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

static bool g_inited = false;
static ssize_t g_client_socket_fd = 0;

/*-------------------------------------------内部接口定义------------------------------------------------*/

static char *base64_encode(const char *input, const size_t length, char *output)  
{  
    if (input == NULL || length < 1) return output;  
  
	*output = '\0';
  
    char *p = (char*)input;  
    char *p_dst = (char*)output;;  
    char *p_end = (char*)input + length;  
    int  loop_count = 0;  
  
    // 0x30 -> 00110000  
    // 0x3C -> 00111100  
    // 0x3F -> 00111111  
    while (p_end - p >= 3) {  
        *p_dst++ = base64_index[( p[0] >> 2 )];  
        *p_dst++ = base64_index[( (p[0] << 4) & 0x30 ) | ( p[1] >> 4 )];  
        *p_dst++ = base64_index[( (p[1] << 2) & 0x3C ) | ( p[2] >> 6 )];  
        *p_dst++ = base64_index[p[2] & 0x3F];  
        p += 3;  
    }  
  
    if (p_end - p > 0) {  
        *p_dst++ = base64_index[(p[0] >> 2)];  
        if (p_end - p == 2) {  
            *p_dst++ = base64_index[( (p[0] << 4) & 0x30 ) | ( p[1] >> 4 )];  
            *p_dst++ = base64_index[(p[1] << 2) & 0x3C];   
            *p_dst++ = '=';  
        } else if (p_end - p == 1) {  
            *p_dst++ = base64_index[(p[1] << 4) & 0x30];  
            *p_dst++ = '=';  
            *p_dst++ = '=';  
        }  
    }  
  
    *p_dst = '\0';  
    return output;  
}  

static char *base64_decode(const char* input, char *output)  
{  
    
    if (input == NULL || output == NULL)   
        return output;  
  
	output[0] = '\0';  
  
    int input_len = strlen(input);  
    if (input_len < 4 || input_len % 4 != 0)   
        return output;  
  
    // 0xFC -> 11111100  
    // 0x03 -> 00000011  
    // 0xF0 -> 11110000  
    // 0x0F -> 00001111  
    // 0xC0 -> 11000000  
    char *p = (char*)input;  
    char *p_out = output;  
    char *p_end = (char*)input + input_len;  
    for (; p < p_end; p += 4) {  
        *p_out++ = ((base64_decode_map[p[0]] << 2) & 0xFC) | ((base64_decode_map[p[1]] >> 4) & 0x03);  
        *p_out++ = ((base64_decode_map[p[1]] << 4) & 0xF0) | ((base64_decode_map[p[2]] >> 2) & 0x0F);  
        *p_out++ = ((base64_decode_map[p[2]] << 6) & 0xC0) | (base64_decode_map[p[3]]);   
    }  
  
    if (*(input + input_len - 2) == '=') {  
        *(p_out - 2) = '\0';  
    } else if (*(input + input_len - 1) == '=') {  
        *(p_out - 1) = '\0';  
    }  
  
    return output;  
} 


static bool client_socket_init(ssize_t* socket_fd, char *ip, uint32_t socket_port)
{
	struct sockaddr_in msg_addr ;
	struct sockaddr_in server_addr;

	*socket_fd= socket(AF_INET, SOCK_STREAM, 0) ; //TCP方式
	if( *socket_fd<0)
	{
		return false;
	} 
	printf("Client client socket init 1111\n");
	memset( &msg_addr, 0,sizeof(msg_addr) );
	msg_addr.sin_family = AF_INET ;
	msg_addr.sin_port = htons( 0 ) ;
	inet_aton(ip, &msg_addr.sin_addr ) ;
	if( bind(*socket_fd, (struct sockaddr *)&msg_addr, sizeof(msg_addr) ) )
	{
		printf("Client client bind error  %s\n", strerror(errno));
		return false;
	}
	
	printf("Client client socket init 2222\n");
	
	memset( &server_addr, 0,sizeof(server_addr) );
	server_addr.sin_family = AF_INET ;
	inet_aton(ip, &server_addr.sin_addr ) ;
	server_addr.sin_port = htons(socket_port) ;
	socklen_t server_addr_len = sizeof(server_addr) ;
	//链接服务端    
	if(connect(*socket_fd,(struct sockaddr *)&server_addr,server_addr_len))
	{
		printf("Client client connect error  %s\n", strerror(errno));
		return false;
	}
	printf("Client client connect ok  %s\n", strerror(errno));
	return true;
}

static bool client_init(char *ip, uint32_t port)
{
	int32_t flags ;
  	bool ret ;
    
    if(!client_socket_init(&g_client_socket_fd, ip, port))
	{
		return false;
	}
	
    //设置是否阻塞 。
    if(SOCKET_IS_BLOCKING)
    {
       flags = fcntl(g_client_socket_fd, F_GETFL, 0); //获取建立的socket的状态 
       fcntl(g_client_socket_fd, F_SETFL, flags&(~O_NONBLOCK));//设置为阻塞状态 
    }
    else
    {
        flags = fcntl(g_client_socket_fd, F_GETFL, 0); //获取建立的socket的状态 
        fcntl(g_client_socket_fd, F_SETFL, flags|O_NONBLOCK);//设置为非阻塞状态
    }
	
	g_inited = true;
	
    return true ;
}

static void client_socket_term(void)
{
	if(g_inited)
	{
		close(g_client_socket_fd) ;
		g_inited = false;
	}
}

static int client_socket_send(char* data, int len)
{
	if(!g_inited)
	{
		return -1;
	}

	if(g_client_socket_fd <= 0)
	{
		return -1;
	}
	
	int ret = 0;
	int writelen = 0;
	char *pData = data;
	int needSendLen = len;
	do{
		ret = send(g_client_socket_fd, pData, needSendLen, 0);
		printf("Client Command send ret = %d sizeof(Socket_command_st) = %d %s\n",ret, needSendLen, strerror(errno));
		pData += ret;
		needSendLen -= ret;
		writelen += ret;
	}while(writelen < len && ret > 0);
	
	printf("Client Command send len = %d writelen = %d\n", len, writelen);
	
	return (writelen > 0 ? writelen : -1);
}

static int client_socket_receive(char* data, int len)
{
	if(!g_inited || g_client_socket_fd <= 0)
	{
		printf("Client matrix_root_command_receive_ex input param error g_inited = %d, g_client_socket_fd = %d\n", g_inited, g_client_socket_fd);
		return -1;
	}
	
	int ret = recv(g_client_socket_fd, data, len, 0);
	printf("Client Command recv ret = %d len = %d %s\n",ret, len, strerror(errno));
	
	return ret;
}


static void* client_socket_heart_beating(void* args) {
	printf("client_socket_heart_beating in\n");
	while (true) {
		char beat[64];
		memset(beat, 0, 64);
		base64_encode("HeartBeating", 12, beat);
		client_socket_send(beat, strlen(beat));
		sleep(5);
	}
	printf("client_socket_heart_beating out\n");
	return NULL;
}

/*-----------------------------------------------------------------------------------------------------*/





int main(void)
{
	char ip[32];
	char port[32];
	memset(ip, 0, 32);
	memset(port, 0, 32);
	printf("\ninput server's ip: ");
	fgets(ip, 32, stdin);
	ip[strlen(ip)-1] = 0;
	printf("\ninput server's port: ");
	fgets(port, 32, stdin);
	port[strlen(port)-1] = 0;
	int sport = atoi(port);
	printf("\ninput server's ip = %s, port = %d\n", ip, sport);
	if (sport <=0) {
		return 0;
	}
	
	client_init(ip, sport);
	
	if(!g_inited || g_client_socket_fd <= 0) {
		return 0;
	}

	static pthread_t thread_handle;
	if (pthread_create(&thread_handle, NULL, client_socket_heart_beating, NULL) != 0)
	{
		printf("client pthread_create failed\n");
	}
	
	while (true) {
		char temp[1024];
		char buff[1024];
		memset(temp, 0, 1024);
		memset(buff, 0, 1024);
		int ret = client_socket_receive(temp, 1024);
		if (ret <= 0) {
			continue;
		}
		
		base64_decode(temp, buff);
		printf("receive data = %s\n", buff);
	}
	
	pthread_join(thread_handle, 0);

    return 0;
}
