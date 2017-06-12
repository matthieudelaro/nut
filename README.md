[![Build Status](https://travis-ci.org/matthieudelaro/nut.svg?branch=master)](https://travis-ci.org/matthieudelaro/nut)

### Whetting Your Appetite
Tired of hearing: "works on my machine"?
Ever experienced headache to install libraries and dependencies?
Ever had to deal with two incompatible versions of a program at once?
Ever wished to try out a new language first, and install it only if it pleases you?
Ever wished to develop for Linux when you use Mac OS or Windows?
Ever wished to develop in Go from the folder of your choice?
Ever wished to have a unified development tool, across all platforms, customizable to any languages?
Ever wished to simplify and share your research on neural networks in Docker, running on GPU, with reproducible results?

### Nut
**Nut** is a command line tool which offers a solution to common frustrations of developers. It hides the complexity of development environments, and extends them with customizable macros. Whether you develop in Swift, Go, Java, or C++, what you need is build/run/test the app. So just do it:
```bash
$ nut --init # create nut.yml file (equivalent of package.json for npm)
$ nut build
$ nut run
$ nut test
```

**Nut** mounts the current folder in a [Docker](https://www.docker.com/) container, and executes commands on your behalf, according to the project configuration. The configuration is read from `nut.yml` file, in the current/parent folder. You can choose the Docker image to use, declare volumes to mount, and define commands (called macros) such as *build*, *run*, and *test*. Nut also [synchronizes timezone](https://github.com/matthieudelaro/nut/wiki/Nut:-Under-The-Hood#synchronize-timezones)

Nut is in early stage of development. It has been tested on Ubuntu, on MacOS with *Docker for Mac* and *Docker Toolbox*, and on Windows with *Docker Toolbox*. Feedbacks and contributions to add features and to [improve Windows support](https://github.com/matthieudelaro/nut/issues/4) are welcome.

Check the [wiki](https://github.com/matthieudelaro/nut/wiki) to learn about Nut [implementation](https://github.com/matthieudelaro/nut/wiki/Nut:-Under-The-Hood), and to read some tutorials (GPU support, Caffe, TensorFlow, etc).

### Share and reuse environments
You can initialize **Nut** with an environment from a GitHub repository:
```bash
$ nut --init --github=matthieudelaro/nutfiles/golang1.6
```
This creates `nut.yml` file that inherites the configuration defined in the nut file at the root of the repository.
This configuration can be overloaded by defining/redefining docker image, macros, mounting points, ... It makes it easy for developers to use libraries and development tools that provide a nut file.

To inspect an environment, you can use `--exec` flag:
```bash
$ # --exec="command to run in the container"
$ nut --exec="pwd"  # will display the path in working directory of the container
$ nut --exec="ls"  # will display the files in the container working directory
$ nut --exec="echo hello && echo world!" --logs  # --logs will display the logs for developers
```
`--exec` flag can be really handy to build and test **Nut** on OSX:
```bash
$ nut build-osx && nut --exec="./nut --init --logs && ls -lah .nut"
```
### Getting Nut
#### Compiling from source
Provided that you use Docker, you don't need to install anything on your computer.
Not even Go!
```bash
# 1 - Download sources and dependencies
    # unix
    docker run --rm -v $PWD:/go/src/github.com/matthieudelaro/ -w /go/src/github.com/matthieudelaro/ matthieudelaro/golang:1.7-cross bash -c 'git clone https://github.com/matthieudelaro/nut.git  --progress && cd nut && govendor sync'

    # windows
    docker run --rm -v ${PWD}:/go/src/github.com/matthieudelaro/ -w /go/src/github.com/matthieudelaro/ matthieudelaro/golang:1.7-cross bash -c 'git clone https://github.com/matthieudelaro/nut.git  --progress && cd nut && govendor sync'

# 2 - Move to nut folder
    cd nut

# 3 - Build Nut
    # Build Nut for Linux, in a container
    docker run -i -t --rm -v $PWD:/go/src/github.com/matthieudelaro/nut -w /go/src/github.com/matthieudelaro/nut matthieudelaro/golang:1.7-cross go build -o nut

    # Build Nut for OSX, in a container
    docker run -i -t --rm -v $PWD:/go/src/github.com/matthieudelaro/nut -w /go/src/github.com/matthieudelaro/nut matthieudelaro/golang:1.7-cross env GOOS=darwin GOARCH=amd64 go build -o nut

    # Build Nut for Windows, in a container.
    docker run -i -t --rm -v ${PWD}:/go/src/github.com/matthieudelaro/nut -w /go/src/github.com/matthieudelaro/nut matthieudelaro/golang:1.7-cross env GOOS=windows GOARCH=amd64 go build -o nut.exe

# Run nut
./nut   # or .\nut.exe on Windows

# Try out Nut
./nut test # will compile and run the tests in a container, according to nut.yml

# Add nut to your PATH
    # Copy it in the path
    sudo cp nut /usr/local/bin/nut # on linux and osx

    # Or modify the path
    echo "PATH=`pwd`:\$PATH" >> ~/.bashrc  # on linux
    echo "PATH=`pwd`:\$PATH" >> ~/.bash_profile  # on osx
```

#### Using NPM
[@RnbWd](https://github.com/RnbWd) developed a npm package with Nut binaries: https://github.com/RnbWd/nut-bin

#### Download Binaries
Manually built binaries for Linux, OSX, and Windows are available in [release](https://github.com/matthieudelaro/nut/tree/manualbuild/release) folder. _It is a temporary solution._


### Nut File Syntax
#### Example
Here is an example of `nut.yml` to develop in Go. You can generate a sample configuration with:

`nut --init`
```yaml
# nut.yml
project_name: nut
enable_gui: true # forward X11 to run graphical application from within the container
                 # On OSX, you have to install an X11 server first : XQuartz (http://www.xquartz.org/) (and you may need to restart your terminal or to reboot, in order to initialize environment variables properly)
                 # On Ubuntu, depending on your config, you may need to run "xhost+" before running nut.

based_on: # configuration can be inherited from:
  github: matthieudelaro/nutfiles/golang1.6 # a GitHub repository
  nut_file_path: ../go1.5/nut.yml # a local file
  # You can inherite either from GitHub or from a file, not both.

docker_image: golang:1.6 # the Docker image. Will override image inherited from file or from GitHub

volumes: # declare folders to mount in the container
  main: # give each folder any name that you like
    host_path: .               # this folder (from your computer) will be mounted as
    container_path: /go/src/project # this folder (in the container)
  shared:
    volume_name: somevolume  # this docker volumes will be mounted as
    container_path: /tmp/shared # this folder (in the container)

macros: # macros define operations that Nut can perform
  build: # call this macro with "nut build"
    usage: build the project
    actions:  # a list of commands to run in the container
    - go build -o nut
    - echo Done
  run: # call this macro with "nut run"
    usage: run the project in the container
    # settings can be overridden for each macro
    enable_current_user: true # login as host user, instead of root
    actions:
    - ./nut
  test:
    usage: test the project
    actions:
    - go test

container_working_directory: /go/src/project # where macros will be executed
work_in_project_folder_as: /go/src/github.com/matthieudelaro/nut
  # Mount the folder of the project (where nut.yml is located) to the given
  # location inside the container. Also set the working directory to this
  # location. This is equivalent to container_working_directory + volume
syntax_version: "7" # Nut evolves quickly ; its configuration file syntax as well.
                    # So nut files are versioned to ensure backward compatibility.

# extra configuration:
privileged: true # run container with --privileged flag
environment: # set environment variables
  var_A: hello # equivalent to: export var_A=hello
  var_B: world
  # environment variables beginning with NUT_ should be reserved for internal usage.
ports: # open ports
  - "3000:3000"
  - 100:100
net: host # docker run --net
uts: host # docker run --uts
security_opts: # docker run --security-opt
  - seccomp=unconfined
devices: # expose devices to the container
  # On OSX and Windows, docker runs into a VM which does not support devices.
  # So Nut supports devices on linux only.
  first:
    host_path: "/dev/1"
    container_path: "/dev/1"
    options: "rw"

```

Here are other instructive [examples](https://github.com/matthieudelaro/nut/blob/master/examples/):
- [Dynamic folder name](https://github.com/matthieudelaro/nutfile_go1.5/blob/master/nut.yml)
- [GUI application](https://github.com/matthieudelaro/nut/blob/master/examples/geary/nut.yml)


#### Guidelines
Nut aims to unify development tools, not to replace compilers.
Nut aims to unify development processes, not to expose language specific requirements.

So, when creating a `nut.yml` file, one should standard names for macros, such as:
- build
- run
- test
- debug
- deploy

As opposed to:
- javac (should be generalized with *build*)
- make (duplicate of *build*)
- do (hum... *Do* what?)
This will keep Nut easy to integrate in text editors and IDEs.

### Support for [nvidia-docker](https://github.com/NVIDIA/nvidia-docker.git)
On Linux, Nut can leverage Nvidia GPUs for your environments. This is useful to use and develop deep learning frameworks, or even to run video games. Due to limitations of Docker on OSX and Windows, Nut does not support GPUs on those platforms.

GPU support relies on [nvidia-docker-plugin](https://github.com/NVIDIA/nvidia-docker/wiki/Using-nvidia-docker-plugin). If it is not running automatically on your machine after [installation](https://github.com/NVIDIA/nvidia-docker/wiki/Installation), you can run it [this way](https://github.com/NVIDIA/nvidia-docker/wiki/Using-nvidia-docker-plugin#usage):
```bash
# Add a system user nvidia-docker
adduser --system --home /var/lib/nvidia-docker nvidia-docker
# Register the plugin with the Docker daemon
mkdir -p /etc/docker/plugins
echo "unix:///var/lib/nvidia-docker/nvidia-docker.sock" > /etc/docker/plugins/nvidia-docker.spec
# Set the mandatory permission
setcap cap_fowner+pe /usr/bin/nvidia-docker-plugin

# Run nvidia-docker-plugin as the nvidia-docker user
sudo -u nvidia-docker nvidia-docker-plugin -s /var/lib/nvidia-docker
```

nvidia-docker-plugin **MUST** be running when you call **Nut**. You can check with:
```bash
curl -s http://0.0.0.0:3476/v1.0/gpu/info  # query the REST API exposed by nvidia-docker-plugin

# should display something like
Driver version:          352.63
Supported CUDA version:  7.5

Device #0
  Model:         GeForce GTX TITAN X
  UUID:          GPU-7e7b6b05-764c-8e74-d867-9a87868d5a1f
  Path:          /dev/nvidia0
  Family:        Maxwell
  Arch:          5.2
  Cores:         3072
  Power:         250 W
  CPU Affinity:  NUMA node0
  PCI
    Bus ID:     0000:01:00.0
    BAR1:       256 MiB
    Bandwidth:  15760 MB/s
  Memory
    ECC:        false
    Global:     12287 MiB
    Constant:   64 KiB
    Shared:     96 KiB
    L2 Cache:   3072 KiB
    Bandwidth:  336480 MB/s
  Clocks
    Cores:        1391 MHz
    Memory:       3505 MHz
  P2P Available:  None
```


### What the Nut???
- build [Nut](https://github.com/matthieudelaro/nut/blob/master/nut.yml) within Nut (I never installed Go, and I'm never going to :)
- build [Docker](https://github.com/matthieudelaro/nut/blob/master/examples/docker/nut.yml)
- build and run [Caffe](https://github.com/matthieudelaro/nut/blob/master/examples/caffe/nut.yml) with `nut build`, `nut test`, `nut train-mnist`.
- compile CUDA code on a Mac Book Air, which hasn't got any Nvidia GPU. Just `nut build`
- test code in a whole infrastructure, by defining a macro running *docker-compose* in a container.
- run linux graphical applications such as [geary](https://github.com/matthieudelaro/nut/blob/master/examples/geary/nut.yml) and [chrome](https://github.com/matthieudelaro/nut/blob/master/examples/chrome/nut.yml) on your Mac after installing [XQuartz](http://www.xquartz.org/):

![Geary application on your Mac](https://camo.githubusercontent.com/b32c086f7da89f3365062f9a6a49b7f64377cb35/687474703a2f2f692e696d6775722e636f6d2f4b6650676d72322e676966)
![Chrome application on your Mac](https://pbs.twimg.com/media/Cl90rCuVYAATcsz.jpg:large)

### Milestones
- create container only once, and store its ID in .nut file
- improve support for Windows
- plugin for Sublime Text, to call `nut run`, `nut build`, and `nut test` from the editor
- create a registery for `nut.yml` files
- see [issues](https://github.com/matthieudelaro/nut/issues)

### Stay Tuned
Wanna receive updates? Or share your thoughts? You can post an issue or follow me on [Twitter](https://twitter.com/matthieudelaro).

### Authors and Contributors
@matthieudelaro and @gdevillele, as well as authors of PRs
