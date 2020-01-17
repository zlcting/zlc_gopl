# **前置条件**
Docker: 要想使用Docker来启动kafka，开发环境提前装好Docker是必须的，我一般在Ubuntu虚拟机上进行开发测试
Docker Compose: kafka依赖zookeeper，使用docker-compose来管理容器依赖
# **Docker镜像**
要想使用Docker安装Kafka，第一件事当然是去Docker hub上找镜像以及使用方法啦。发现kafka并不像mysql或者redis那样有官方镜像，不过Google一下后发现可以选择知名的三方镜像wurstmeister/kafka

wurstmeister/kafka在Github上更新还算频繁，目前使用kafka版本是1.1.0

# **安装**
1.参考官方测试用的docker-compose.yml直接在自定义的目录位置新建docker-compose的配置文件
```
touch ~/docker/kafka/docker-compose.yml
```
```
version: '2.1'
services:
  zookeeper:
    image: wurstmeister/zookeeper
    ports:
      - "2181"
  kafka:
    image: wurstmeister/kafka
    ports:
      - "9092"
    environment:
      KAFKA_ADVERTISED_HOST_NAME: 192.168.5.139
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock

```
注意： KAFKA_ADVERTISED_HOST_NAME 需要配置为宿主机的ip

2.docker-compose 启动kafka
```
root@ubuntu:~/docker/kafka# docker-compose up -d
```
启动完之后通过docker ps可以看到启动了一个zookeeper容器和一个kafka容器

3启动多个kafka节点，比如3
```
root@ubuntu:~/docker/kafka# docker-compose scale kafka=3
```
# **验证**
1.首先进入到一个kafka容器中，例如: kafka_kafka_1
```
root@ubuntu:~/docker/kafka# docker exec -it kafka_kafka_1 /bin/bash
```
2创建一个topic并查看，需要指定zookeeper的容器名(这里是kafka_zookeeper_1)，topic的名字为test

```
$KAFKA_HOME/bin/kafka-topics.sh --create --topic test --zookeeper kafka_zookeeper_1:2181 --replication-factor 1 --partitions 1

$KAFKA_HOME/bin/kafka-topics.sh --list --zookeeper kafka_zookeeper_1:2181
```
3.发布消息，输入几条消息后，按^C退出发布
```
$KAFKA_HOME/bin/kafka-console-producer.sh --topic=test --broker-list kafka_kafka_1:9092
```
4.接受消息
```
$KAFKA_HOME/bin/kafka-console-consumer.sh --bootstrap-server kafka_kafka_1:9092 --from-beginning --topic test
```
如果接收到了发布的消息，那么说明部署正常，可以正式使用了。


docker-compose部署zk集群、kafka集群以及kafka-manager

docker-compose的详细使用方法参考下面这个博客
https://blog.51cto.com/9291927/2310444

```xml
version: '2.1'
services:
  zookeeper:
    image: wurstmeister/zookeeper
    ports:
      - "2181:2181"
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
  kafka:
    image: wurstmeister/kafka
    ports:
      - "9092:9092"
    depends_on:
      - "zookeeper"
    environment:
      KAFKA_BROKER_NO: 0
      KAFKA_ADVERTISED_HOST_NAME: 10.208.86.17
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://10.208.86.17:9092
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_PORT: 9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```      