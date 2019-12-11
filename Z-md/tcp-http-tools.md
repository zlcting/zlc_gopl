# 网络协议应用浅析以及网络工具使用

TCP 是一个可靠的（reliable）、面向连接的（connection-oriented）、基于字节流（byte-stream）、全双工的（full-duplex）协议

我们熟悉或者经常见到的网络协议，不分门别类的大致罗列一下有，tcp ip http udp dns smtp ftp 
其中一tcp ip http 最为常见


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

打开 wireshark，在弹出的选项中，选中 en0 网卡，在过滤器中输入host www.baidu.com，只抓取与百度服务器通信的数据包。
![图片](http://raw.githubusercontent.com/zlcting/zlc_gopl/master/Z-md/img/net-tcp-tool-fenceng.png)

Ethernet II：网络接口层以太网帧头部信息
Internet Protocol Version 4：互联网层 IP 包头部信息
Transmission Control Protocol：传输层的数据段头部信息，此处是 TCP 协议
Hypertext Transfer Protocol：应用层 HTTP 的信息



