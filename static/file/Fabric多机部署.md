## 一.环境准备

我们现在部署的是2Peer+3Orderer的fabric网络组织架构。为此我们需要准备5台机器。我们可以开5台虚拟机（可以在VMware14或VirtualBox工具下创建安装5个Ubuntu18.04系统的虚拟机），也可以购买5台云服务器（阿里云服务器即可或者其他），不管怎么样，我们需要这5台机器网络能够互通，而且安装相同的系统，我们用的是Ubuntu 18.04版。为了方便，我建议先启用1台虚拟机，在其中把准备工作做完，然后基于这台虚拟机，再复制出4台即可。这里是我用到5台Server的主机名（角色）和IP：

| orderer0.bss.com   | 192.168.1.221 |
| ------------------ | :------------ |
| orderer1.bss.com   | 192.168.1.222 |
| orderer2.bss.com   | 192.168.1.223 |
| peer0.org1.bss.com | 192.168.1.224 |
| peer0.org2.bss.com | 192.168.1.225 |

接下来我们需要准备软件环境（每台机器都要安装一遍）

### 1.配置ubuntu系统

不管你是用VirtualBox或VMware工具创建安装ubuntu18.04系统的虚拟机，还是直接用云服务器，你最好先确定一下每一台机器ubuntu18.04系统的apt source得是国内的比较好，不然如果是国外的话会很慢很慢的。具体的做法是

`sudo vi /etc/apt/sources.list`

打开这个apt源列表，如果其中看到是http://us.xxxxxx之类的，那么就是外国的，如果看到是[http://cn.xxxxx](http://cn.xxxxx/)之类的，那么就不用换的。我的是美国的源，所以需要做一下批量的替换。在命令模式下，输入：

<code>:%s/us./cn./g</code>

就可以把所有的us.改为cn.了。然后输入:wq!即可保存退出。

<code>sudo apt-get update</code>

更新一下源。

然后安装ssh，这样接下来就可以用putty或者SecureCRT之类的客户端远程连接Ubuntu了。

<code>sudo apt-get install ssh</code>

### 2.安装go

Ubuntu的apt-get虽然提供了Go的安装，但是版本比较旧，最好的方法还是参考go语言中文网https://studygolang.com/dl ，下载所需要的Go版本。在这里bss1.4.1需要go的版本在1.11.5及以上版本。具体涉及到的命令包括：

<code>wget https://studygolang.com/dl/golang/go1.12.5.linux-amd64.tar.gz</code>

或者

<code>sudo curl -O https://storage.googleapis.com/golang/go1.12.5.linux-amd64.tar.gz</code>

<code>sudo tar -C /~ -xzf go1.12.5.linux-amd64.tar.gz</code>

注：不要使用apt方式安装go，apt的go版本太低了！

解压后，生成一个go目录。

用命令行mv将该目录移到目录/usr/local下：

<code>mv /~/go /usr/local</code>

接下来编辑当前用户的环境变量,首先创建相应的目录：

<code>cd ~<br>mkdir -p /~/gocode<br>cd /~/gocode<br>mkdir -p src<br>mkdir -p bin<br>mkdir -p pkg<br>vim ~/.profile</code>

添加以下内容：

<code>#根目录<br>export GOROOT=/usr/local/go<br>#bin目录<br>export GOBIN=$GOROOT/bin<br>#工作目录<br>export GOPATH=/~/gocode<br>export PATH=$PATH:$GOPATH:$GOBIN:$GOPATH<br></code>或者
<code>export GOROOT=/usr/local/go<br>export GOBIN=$GOPATH/bin<br>export GOPATH=$HOME/gocode<br>export PATH=$PATH:$GOROOT/bin<br>export PATH=$PATH:$GOPATH/bin<br>export PATH=$PATH:$GOPATH</code>

编辑保存并退出vi后，记得把这些环境载入：

<code>source /etc/profile</code>

运行以下命令查看当前go的版本，如果能够显示go版本，那么说明我们的go安装成功.

<code>go version<br>go version go1.12.5 linux/amd64</code>
### 3.安装docker

#### 1.1使用DaoClound设置docker镜像

安装Docker也会遇到外国网络慢的问题，幸好国内有很好的镜像，推荐DaoClound，安装Docker的命令是：

<code>curl -sSL https://get.daocloud.io/docker | sh</code>

安装完成后，修改当前用户（使用的用户叫fabric）权限，运行以下脚本将当前用户添加到Docker的组中

<code>sudo usermod -aG docker fabric</code>

(fabric是你的ubuntu系统用户名)
注销并重新登录当前用户，接下来修改 Docker 服务配置（/etc/default/docker 文件，如果没有这个文件用vi命令也会创建该文件的）

<code>sudo vi /etc/default/docker</code>

添加以下内容：
<code>DOCKER_OPTS="$DOCKER_OPTS -H tcp://0.0.0.0:2375 -H unix:///var/run/docker.sock --api-cors-header='*'"</code>

接下来就需要设置国内的Docker镜像地址，需要注册一个账号，然后在加速器页面提供了设置Docker镜像的脚本，加速器页面是：

<code>https://www.daocloud.io/mirror </code>

我提供的脚本是：

<code>curl -sSL https://get.daocloud.io/daotools/set_mirror.sh | sh -s http://d4cc5789.m.daocloud.io</code>
运行完脚本后，重启Docker服务
<code>sudo service docker restart</code>