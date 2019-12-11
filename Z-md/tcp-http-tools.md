# 网络协议应用浅析以及网络工具使用

![图片](https://user-gold-cdn.xitu.io/2019/5/20/16ad5182f90d0bb5?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)



![图片](https://user-gold-cdn.xitu.io/2019/10/11/16dbb2fbdaebd148?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

收到 IP 数据包解析以后，它怎么知道这个分组应该投递到上层的哪一个协议（UDP 或 TCP）?

网络访问层中有个Protocol 标示 是TCP


#####三次握手
tcp标志位，有6种标示：
SYN(synchronous建立联机) 
ACK(acknowledgement 确认)
RST(reset重置)  
PSH(push传送) 
FIN(finish结束) 
URG(urgent紧急)
Sequence number(顺序号码) 
Acknowledge number(确认号码)

![图片](https://user-gold-cdn.xitu.io/2019/3/16/16985bd53967c3b2?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

第一次握手：主机A发送位码为syn＝1，随机产生seq number=1234567的数据包到服务器，主机B由SYN=1知道，A要求建立联机；

第二次握手：主机B收到请求后要确认联机信息，向A发送ack number=(主机A的seq+1)，syn=1，ack=1，随机产生seq=7654321的包；

第三次握手：主机A收到后检查ack number是否正确，即第一次发送的seq number+1，以及位码ack是否为1，若正确，主机A会再发送ack number=(主机B的seq+1)，ack=1，主机B收到后确认seq值与ack=1则连接建立成功。


#####TCP 是一个可靠的（reliable）、面向连接的（connection-oriented）、基于字节流（byte-stream）、全双工的（full-duplex）协议

####TCP 协议是可靠的
P 是一种无连接、不可靠的协议：它尽最大可能将数据报从发送者传输给接收者，但并不保证包到达的顺序会与它们被传输的顺序一致，也不保证包是否重复，甚至都不保证包是否会达到接收者。

TCP 要想在 IP 基础上构建可靠的传输层协议，必须有一个复杂的机制来保障可靠性。 主要有下面几个方面：

对每个包提供校验和
包的序列号解决了接收数据的乱序、重复问题
超时重传
流量控制、拥塞控制
校验和（checksum） 每个 TCP 包首部中都有两字节用来表示校验和，防止在传输过程中有损坏。如果收到一个校验和有差错的报文，TCP 不会发送任何确认直接丢弃它，等待发送端重传。
![图片](https://user-gold-cdn.xitu.io/2019/10/9/16dafd4097c7d058?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

######包的序列号保证了接收数据的乱序和重复问题 
假设我们往 TCP 套接字里写 3000 字节的数据导致 TCP发送了 3 个数据包，每个数据包大小为 1000 字节：第一个包序列号为[1~1001)，第二个包序列号为 [1001~2001)，第三个包序号为[2001~3001)
![图片](https://user-gold-cdn.xitu.io/2019/3/16/16985bd5397b180a?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)
假如因为网络的原因导致第二个、第三个包先到接收端，第一个包最后才到，接收端也不会因为他们到达的顺序不一致把包弄错，TCP 会根据他们的序号进行重新的排列然后把结果传递给上层应用程序。

如果 TCP 接收到重复的数据，可能的原因是超时重传了两次但这个包并没有丢失，接收端会收到两次同样的数据，它能够根据包序号丢弃重复的数据。



#####0x03 TCP 是面向字节流的协议
TCP 是一种字节流（byte-stream）协议，流的含义是没有固定的报文边界。

假设你调用 2 次 write 函数往 socket 里依次写 500 字节、800 字节。write 函数只是把字节拷贝到内核缓冲区，最终会以多少条报文发送出去是不确定的，如下图所示
![图片](https://user-gold-cdn.xitu.io/2019/3/17/1698a074292fb212?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

情况 1：分为两条报文依次发出去 500 字节 和 800 字节数据，也有
情况 2：两部分数据合并为一个长度为 1300 字节的报文，一次发送
情况 3：第一部分的 500 字节与第二部分的 500 字节合并为一个长度为 1000 字节的报文，第二部分剩下的 300 字节单独作为一个报文发送
情况 4：第一部分的 400 字节单独发送，剩下100字节与第二部分的 800 字节合并为一个 900 字节的包一起发送。

上面出现的情况取决于诸多因素：路径最大传输单元 MTU、发送窗口大小、拥塞窗口大小等。

当接收方从 TCP 套接字读数据时，它是没法得知对方每次写入的字节是多少的。接收端可能分2 次每次 650 字节读取，也有可能先分三次，一次 100 字节，一次 200 字节，一次 1000 字节进行读取。
上面出现的情况取决于诸多因素：路径最大传输单元 MTU、发送窗口大小、拥塞窗口大小等。

####0x04 TCP 是全双工的协议
在 TCP 中发送端和接收端可以是客户端/服务端，也可以是服务器/客户端，通信的双方在任意时刻既可以是接收数据也可以是发送数据，每个方向的数据流都独立管理序列号、滑动窗口大小、MSS 等信息。


###数据包大小对网络的影响——MTU与MSS的奥秘


### 第一节 兄弟你还好吗？ 协议中的  keepalive 

我们先来可一个问题，我们浏览器打开一个页面的时候，到底打开了多少个tcp链接呢？


##### 永远记住 TCP 不是轮询的协议

TCP 的 keepalive 与 HTTP 的 keep-alive 有什么区别？

###socket协议 

### 网络工具使用介绍
#####一.命令行工具
命令一：telnet
命令二：netcat
nc -l -p 8080
nc ip 8080
命令三：netstat

#####二.TCPDump 基础 

tcpdump介绍一个案例

#####三.Wireshark 使用实战

TCP/IP 网络分层
记得在学习计算机网络课程的时候，一上来就开始讲分层模型了，当时死记硬背的各个层的名字很快就忘光了，不明白到底分层有什么用。纵观计算机和分布式系统，你会发现「计算机的问题都可以通过增加一个虚拟层来解决，如果不行，那就两个」

下面用 wireshark 抓包的方式来开始看网络分层。
![图片](https://user-gold-cdn.xitu.io/2019/5/20/16ad5181c0a6eb2c?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

打开 wireshark，在弹出的选项中，选中 en0 网卡，在过滤器中输入host www.baidu.com，只抓取与百度服务器通信的数据包。

![图片](https://user-gold-cdn.xitu.io/2019/5/20/16ad5181cb911e12?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

Ethernet II：网络接口层以太网帧头部信息
Internet Protocol Version 4：互联网层 IP 包头部信息
Transmission Control Protocol：传输层的数据段头部信息，此处是 TCP 协议
Hypertext Transfer Protocol：应用层 HTTP 的信息

####洪水攻击原理
scrapy工具
https://github.com/secdev/scapy

