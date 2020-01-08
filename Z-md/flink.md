
注：https://github.com/apache/flink  github源码

一、集群部署和启动
启动：
bin/start-cluster.sh
添加JobManager
bin/jobmanager.sh ((start|start-foreground) cluster)|stop|stop-all
添加TaskManager
bin/taskmanager.sh start|start-foreground|stop|stop-all
停止服务
bin/stop-cluster.sh

localhost/

二、开发环境
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

3.pox.xml文件中配置基础参数

```xml
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <groupId>org.example</groupId>
    <artifactId>learnflink</artifactId>
    <version>1.0-SNAPSHOT</version>

    <dependencies>
    <dependency>
        <groupId>org.apache.flink</groupId>
        <artifactId>flink-java</artifactId>
        <version>1.9.1</version>
<!--        <scope>provided</scope>-->
    </dependency>
    <dependency>
        <groupId>org.apache.flink</groupId>
        <artifactId>flink-streaming-java_2.11</artifactId>
        <version>1.9.1</version>
<!--        <scope>provided</scope>-->
    </dependency>
    <dependency>
        <groupId>org.apache.flink</groupId>
        <artifactId>flink-clients_2.11</artifactId>
        <version>1.9.1</version>
<!--        <scope>provided</scope>-->
    </dependency>
    </dependencies>

    <build>
        <plugins>
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-shade-plugin</artifactId>
                <version>3.0.0</version>
                <executions>
                    <execution>
                        <phase>package</phase>
                        <goals>
                            <goal>shade</goal>
                        </goals>
                        <configuration>
                            <artifactSet>
                                <excludes>
                                    <exclude>com.google.code.findbugs:jsr305</exclude>
                                    <exclude>org.slf4j:*</exclude>
                                    <exclude>log4j:*</exclude>
                                </excludes>
                            </artifactSet>
                            <filters>
                                <filter>
                                    <!-- Do not copy the signatures in the META-INF folder.
                                    Otherwise, this might cause SecurityExceptions when using the JAR. -->
                                    <artifact>*:*</artifact>
                                    <excludes>
                                        <exclude>META-INF/*.SF</exclude>
                                        <exclude>META-INF/*.DSA</exclude>
                                        <exclude>META-INF/*.RSA</exclude>
                                    </excludes>
                                </filter>
                            </filters>
                            <transformers>
                                <transformer implementation="org.apache.maven.plugins.shade.resource.ManifestResourceTransformer">
                                    <mainClass>my.programs.main.clazz</mainClass>
                                </transformer>
                            </transformers>
                        </configuration>
                    </execution>
                </executions>
            </plugin>
        </plugins>
    </build>

</project>
```
然后点击idea 的 Import Changes 提示 加载flink包文件

三.应用开发
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




