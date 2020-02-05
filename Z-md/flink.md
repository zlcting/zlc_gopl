
注：https://github.com/apache/flink  github源码
启动tcp服务： nc -l -p 9001

### 一、开发环境
开发工具推荐 idea

1.使用idea创建maven项目


2.安装maven工具 （在maven/conf/settings.xml 修改阿里云镜像可以使加载更快 ）
```xml
    <mirror>  
        <id>nexus-aliyun</id>  
        <mirrorOf>central</mirrorOf>    
        <name>Nexus aliyun</name>  
        <url>http://maven.aliyun.com/nexus/content/groups/public</url>  
    </mirror>

    <localRepository>/usr/local/maven/repository（替换自己maven路径）</localRepository>
```
idea 菜单修改 file>setting>Build...>maven>user setting file >"/usr/local/maven/conf/settings.xml" 

mvn archetype:generate                               \
-DarchetypeGroupId=org.apache.flink              \
-DarchetypeArtifactId=flink-quickstart-java      \
-DarchetypeVersion=1.9.0 \
-DarchetypeVersion=local


然后点击idea 的 Import Changes 提示 加载flink包文件

### 二.应用开发
1.代码步骤
1).获得一个execution environment，
2).加载/创建初始数据，
3).指定此数据的转换，
4).指定放置计算结果的位置，
5).触发程序执行

2.Lazy Evaluation 延迟执行 

execute() 是调用执行代码的入口函数

[3.Specifying Keys 指定keys](https://ci.apache.org/projects/flink/flink-docs-release-1.9/zh/dev/api_concepts.html#specifying-keys
 "指定keys")

关键字：join, coGroup, keyBy, groupBy 

1）Define keys for Tuples 通过元组来指定keys
```java
    text.flatMap(new FlatMapFunction<String, Tuple2<String,Integer>>() {
        @Override
        public void flatMap(String value, Collector<Tuple2<String, Integer>> collector) throws Exception {
            String[] tokens = value.toLowerCase().split("/t");
            for (String token:tokens){
                if (token.length()>0){
                    collector.collect(new Tuple2<String, Integer>(token,1));
                }
            }
        }
    }).keyBy(0).timeWindow(Time.seconds(5)).sum(1).print();
```
2）Define keys using Field Expressions 使用字段表达式定义keys
在下面的例子中，我们有一个WC POJO，它有两个字段“word”和“count”。要按字段单词分组，只需将其名称传递给keyBy（）函数。

例子：
```java
public static void main(String[] args) throws Exception{
        StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();

        DataStreamSource<String> text =  env.socketTextStream("127.0.0.1",9001);

        text.flatMap(new FlatMapFunction<String, WC>() {
            @Override
            public void flatMap(String value, Collector<WC> collector) throws Exception {
                String[] tokens = value.toLowerCase().split("/t");
                for (String token:tokens){
                    if (token.length()>0){
                        collector.collect(new WC(token,1));
                    }
                }
            }
        }).keyBy("word")
                .timeWindow(Time.seconds(5))
                .sum("count")
                .print();

        env.execute("zhixing");
    }
    public static class WC{
        private String word;
        private int count;
        public WC(){}
        public WC(String word,int count){
            this.word=word;
            this.count=count;
        }

        @Override
        public String toString() {
            return "WC{" +
                    "word='" + word + '\'' +
                    ", count=" + count +
                    '}';
        }

        public String getWord() {
            return word;
        }

        public void setWord(String word) {
            this.word = word;
        }

        public int getCount() {
            return count;
        }

        public void setCount(int count) {
            this.count = count;
        }
    }
  ```
3）.Define keys using Key Selector Functions 使用选择函数来指定keys

 ```java
    // some ordinary POJO
    public class WC {public String word; public int count;}
    DataStream<WC> words = // [...]
    KeyedStream<WC> keyed = words
    .keyBy(new KeySelector<WC, String>() {
        public String getKey(WC wc) { return wc.word; }
    });
 ```
[3.Specifying Transformation Functions 指定转换函数](https://ci.apache.org/projects/flink/flink-docs-release-1.9/zh/dev/api_concepts.html#specifying-transformation-functions
 "指定转换函数")

1）Implementing an interface 使用接口方法

 ```java
    DataStreamSource<String> text =  env.socketTextStream("127.0.0.1",9001);

     text.flatMap(new MyFlatMapFunction()).keyBy("word")
                .timeWindow(Time.seconds(5))
                .sum("count")
                .print();

    public static class MyFlatMapFunction implements FlatMapFunction<String, WC> {
        @Override
        public void flatMap(String value, Collector<WC> collector) throws Exception {
            String[] tokens = value.toLowerCase().split("/t");
            for (String token:tokens){
                if (token.length()>0){
                    collector.collect(new WC(token,1));
                }
            }
        }
    }
 ```
2).Rich functions
All transformations that require a user-defined function can instead take as argument a rich function. 
所有需要用户定义函数的转换都可以将富函数作为参数。

例如，代替

 ```java
    class MyMapFunction implements MapFunction<String, Integer> {
    public Integer map(String value) { return Integer.parseInt(value); }
    };
 ```
可以写成
 ```java
    class MyMapFunction extends RichMapFunction<String, Integer> {
    public Integer map(String value) { return Integer.parseInt(value); }
    };
 ```
并且也可以使用匿名函数
 ```java
    data.map(new MyMapFunction());
 ```

 ```java
    data.map (new RichMapFunction<String, Integer>() {
    public Integer map(String value) { return Integer.parseInt(value); }
    });
 ```
Rich functions provide, in addition to the user-defined function (map, reduce, etc), four methods: open, close, getRuntimeContext, and setRuntimeContext. These are useful for parameterizing the function (see Passing Parameters to Functions), creating and finalizing local state, accessing broadcast variables (see Broadcast Variables), and for accessing runtime information such as accumulators and counters (see Accumulators and Counters), and information on iterations (see Iterations).



[4.Flink DataSet API]( https://ci.apache.org/projects/flink/flink-docs-release-1.9/zh/dev/batch/
 "")




算子：
1.Flink单数据流基本转换：map、filter、flatMap
https://link.zhihu.com/?target=https%3A//mp.weixin.qq.com/s/z8L6QU1ZWW1-O2cn8ixcmg
2.Flink基于Key的分组转换：keyBy、reduce和aggregations
https://link.zhihu.com/?target=https%3A//mp.weixin.qq.com/s/2vcKteQIyj31sVrSg1R_2Q
3.Flink多数据流转换：union和connect
https://link.zhihu.com/?target=https%3A//mp.weixin.qq.com/s/vz94e-TAKa1da9Nd8O6kFw
4.Flink并行度和数据重分配
https://link.zhihu.com/?target=https%3A//mp.weixin.qq.com/s/c4vtqbxhqqVq0hg8C2KGvg

https://zhuanlan.zhihu.com/p/100416194


map算子对一个DataStream中的每个元素使用用户自定义的map函数进行处理，每个输入元素对应一个输出元素，最终整个数据流被转换成一个新的DataStream。输出的数据流DataStream[OUT]类型可能和输入的数据流DataStream[IN]不同

flatMap算子和map有些相似，输入都是数据流中的每个元素，与之不同的是，flatMap的输出可以是零个、一个或多个元素，当输出元素是一个列表时，flatMap会将列表展平。

## 也就是说 flatMap是将输入的数据整体进行整理而map是对每一个输入数据进行处理

keyBy算子将DataStream转换成一个KeyedStream。KeyedStream是一种特殊的DataStream，事实上，KeyedStream继承了DataStream，DataStream的各元素随机分布在各Task Slot中，KeyedStream的各元素按照Key分组，分配到各Task Slot中。我们需要向keyBy算子传递一个参数，以告知Flink以什么字段作为Key进行分组。
我们可以使用数字位置来指定Key：

 ```Scala
val dataStream: DataStream[(Int, Double)] = senv.fromElements((1, 1.0), (2, 3.2), (1, 5.5), (3, 10.0), (3, 12.5))
// 使用数字位置定义Key 按照第一个字段进行分组
val keyedStream = dataStream.keyBy(0)
 ```

自定义sink总结
1）RichSinkFunction<T> T就是你想要写入对象的类型
2）重写方法
open/close 生命周期方法
invoke 每条记录执行一次



### 三、编译flink源代码 

 maven的版本需要3.3.x 和 jdk1.8（必要条件）
 在flink目录执行编译程序
 ```maven
    //快速编译
    mvn clean install -DskipTests
    //跳过单元测试
    mvn clean install -DskipTests -Dfast
    //指定Hadoop版本
    mvn clean install -DskipTests -Dhadoop.version=2.6.5
 ```
### 四、安装部署
下载地址
https://flink.apache.org/downloads.html

 ```
    cd ~/Downloads        # Go to download directory
    $ tar xzf flink-*.tgz   # Unpack the downloaded archive
    $ cd flink-1.9.0
 ```
##### 1.单机实例
启动：
bin/start-cluster.sh

打开http://localhost:8081 你会获得一个ui界面

![图片](https://ci.apache.org/projects/flink/flink-docs-release-1.9/page/img/quickstart-setup/jobmanager-1.png)

执行一个案例
第一步启动一个tcp服务
 ```
    nc -l -p 9001
 ```
提交一个flink 程序 jar包
 ```
 $ ./bin/flink run examples/streaming/SocketWindowWordCount.jar --port 9000
Starting execution of program
```
[源码内容地址](https://ci.apache.org/projects/flink/flink-docs-release-1.9/getting-started/tutorials/local_setup.html#read-the-code "")


 ```
 $ nc -l 9000
    lorem ipsum
    ipsum ipsum ipsum
    bye
 ```
输出的日志记录在
 ```
 $ tail -f log/flink-*-taskexecutor-*.out
    lorem : 1
    bye : 1
    ipsum : 4
 ```

![图片](https://ci.apache.org/projects/flink/flink-docs-release-1.9/page/img/quickstart-setup/jobmanager-2.png)

![图片](https://ci.apache.org/projects/flink/flink-docs-release-1.9/page/img/quickstart-setup/jobmanager-3.png)


停止服务
bin/stop-cluster.sh



添加JobManager
bin/jobmanager.sh ((start|start-foreground) cluster)|stop|stop-all



添加TaskManager
bin/taskmanager.sh start|start-foreground|stop|stop-all

JobManager：

TaskManager：


#### 2.分布式搭建之 Standalone Cluster
第一步先按照单机实例安装。

![图片](https://ci.apache.org/projects/flink/flink-docs-release-1.9/page/img/quickstart_cluster.png)

在目录/path/to/flink/conf/
flink-conf.yaml配置
主节点ip
```
jobmanager.rpc.address: 10.0.0.1
```
/path/to/flink/
conf/slaves
```
10.0.0.2
10.0.0.3
```
常用配置
jobmanager.heap.size jobmanager可用的内存 

taskmanager.heap.size taskmanager可用的内存

taskmanager.numberOfTaskSlots 每个机器可用的cpu个数
parallelism.default 任务的并行度

io.tmp.dirs 临时文件的目录


三 Hadoop yarn flink安装
1.前置条件
准备工作需要下载安装Hadoop环境

yarn session 前置条件
1）、at least Apache Hadoop 2.2
2）、HDFS (Hadoop Distributed File System) (or another distributed file system supported by Hadoop)
设置环境变量
3）、export HADOOP_CLASSPATH=`hadoop classpath`

4）、Start a YARN session where the job manager gets 1 GB of heap space and the task managers 4 GB of heap space assigned:


下载对应的Hadoop的flink包
http://flink.apache.org/downloads.html

编译flink教程可以参考https://blog.csdn.net/aikfk/article/details/102494709
https://www.jianshu.com/p/f9356a5a388d

当然也可用直接从官网下载对应Hadoop的版本

2.Run a Flink job on YARN
```
tar xvzf flink-1.9.0-bin-hadoop2.tgz
cd flink-1.9.0/
./bin/flink run -m yarn-cluster -p 4 -yjm 1024m -ytm 4096m ./examples/batch/WordCount.jar
```


3.Start Flink Session
./bin/yarn-session.sh

```
Usage:
   Optional
     -D <arg>                        Dynamic properties
     -d,--detached                   Start detached
     -jm,--jobManagerMemory <arg>    Memory for JobManager Container with optional unit (default: MB)
     -nm,--name                      Set a custom name for the application on YARN
     -at,--applicationType           Set a custom application type on YARN
     -q,--query                      Display available YARN resources (memory, cores)
     -qu,--queue <arg>               Specify YARN queue.
     -s,--slots <arg>                Number of slots per TaskManager
     -tm,--taskManagerMemory <arg>   Memory per TaskManager Container with optional unit (default: MB)
     -z,--zookeeperNamespace <arg>   Namespace to create the Zookeeper sub-paths for HA mode
```
Example
```
wget -O LICENSE-2.0.txt http://www.apache.org/licenses/LICENSE-2.0.txt
hadoop fs -copyFromLocal LICENSE-2.0.txt hdfs:/// ...
./bin/flink run ./examples/batch/WordCount.jar \
       --input hdfs:///..../LICENSE-2.0.txt --output hdfs:///.../wordcount-result.txt
```

yarn session 上执行的taskmanager 开始 使用8G内存  32进程并发
     ```
     ./bin/yarn-session.sh -tm 8192 -s 32
     ```



![图片](https://ci.apache.org/projects/flink/flink-docs-release-1.9/fig/FlinkOnYarn.svg)
