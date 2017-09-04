### Wetting Your Appetite
Tired of hearing: "works on my machine"?  
Ever experienced headache to install libraries and dependencies?  
Ever had to deal with two incompatible versions of a program at once?  
Ever wished to try out a new language first, and install it only if it pleases you?  
Ever wished to develop for Linux when you use Mac OS or Windows?  
Ever wished to develop in Go from the folder of your choice?  
Ever wished to have a unified development tool, across all platforms, customizable to any languages?  

### Nut
**Nut** is a command line tool which offers a solution to common frustrations of developers. It hides the complexity of development environments, and extends them with customizable macros. Whether you develop in Swift, Go, Java, or C++, what you need is build/run/test the app. So just do it:

    $ nut --init # create nut.yml file (equivalent of package.json for npm)
    $ nut build
    $ nut run
    $ nut test

**Nut** mounts the current folder in a [Docker](https://www.docker.com/) container, and executes commands on your behalf, according to the project configuration. The configuration is read from `nut.yml` file, in the current folder. You can choose the Docker image to use, declare volumes to mount, and define commands (called macros) such as *build*, *run*, and *test*.

Nut is in early stage of development. It has been tested on Ubuntu and on MacOS with *Docker for Mac*. Feedbacks and contributions to add features and to make Nut run on other systems are welcome (Windows, Docker Toolbox, etc).

### Share and reuse environments
You can initialize **Nut** with an environment from a GitHub repository:

    $ nut --init --github=matthieudelaro/nutfile_go1.6
    
This creates `nut.yml` file that inherites the configuration defined in the nut file at the root of the repository.
This configuration can be overloaded by defining/redefining docker image, macros, mounting points, ... It makes it easy for developers to use libraries and development tools that provide a nut file.

To inspect an environment, you can use `--exec` flag:
    
    $ # --exec="command to run in the container"
    $ nut --exec="pwd"  # will display the path in working directory of the container 
    $ nut --exec="ls"  # will display the files in the container working directory
    $ nut --exec="echo hello && echo world!" --logs  # --logs will display the logs for developers

`--exec` flag can be really handy to build and test **Nut** on OSX:

    $ nut build-osx && nut --exec="./nut --init --logs && ls -lah .nut"

### Getting Nut
#### Compile from source
Provided that you use Docker, you don't need to install anything on your computer.  
Not even Go!
```bash
# 1 - Download sources
git clone git@github.com:matthieudelaro/nut.git --recursive
    # if you don't have git, you can download a zip file from http://matthieudelaro.github.io/nut/

# 2 - Move to nut folder
cd nut

# 3 - Build Nut
    # Build Nut for Linux, in a container
    docker run -i -t --rm -v $PWD:/go/src/github.com/matthieudelaro/nut -w /go/src/github.com/matthieudelaro/nut golang:1.6 go build -o nut

    # Build Nut for OSX, in a container
    docker run -i -t --rm -v $PWD:/go/src/github.com/matthieudelaro/nut -w /go/src/github.com/matthieudelaro/nut golang:1.6 env GOOS=darwin GOARCH=amd64 go build -o nut

    # Build Nut for Windows, in a container
    docker run -i -t --rm -v $PWD:/go/src/github.com/matthieudelaro/nut -w /go/src/github.com/matthieudelaro/nut golang:1.6 env GOOS=windows GOARCH=amd64 go build -o nut

# Run nut
./nut

# Try out Nut
./nut test # will compile and run the tests in a container, according to nut.yml

# Add nut to your PATH
    # Copy it in the path
    sudo cp nut /usr/local/bin/nut # on linux and osx

    # Or modify the path
    echo "PATH=`pwd`:\$PATH" >> ~/.bashrc  # on linux
    echo "PATH=`pwd`:\$PATH" >> ~/.bash_profile  # on osx
```

### Nut File Syntax
#### Example
Here is an example of `nut.yml` to develop in Go. You can generate a sample configuration with  :

`nut --init`
```yaml
# nut.yml
project_name: nut
enable_gui: yes # forward X11 to run graphical application from within the container
                # On OSX, you have to install an X11 server first : XQuartz (http://www.xquartz.org/) (and you may need to restart your terminal or to reboot, in order to initialize environment variables properly)
                # On Ubuntu, depending on your config, you may need to run "xhost+" before running nut.
privileged: true # run container with --privileged flag

based_on: # configuration can be inherited from:
  github: matthieudelaro/nutfile_go1.6 # a GitHub repository
  nut_file_path: ../go1.5/nut.yml # a local file
  # You can inherite either from GitHub or from a file, not both.
  docker_image: golang:1.6 # a Docker image. Will override image set on GitHub

mount: # declare folders to mount in the container
  main: # give each folder any name that you like
  - .               # this folder (from your computer) will be mounted as
  - /go/src/project # this folder (in the container)

macros: # macros define operations that Nut can perform
  build: # call this macro with "nut build"
    usage: build the project
    actions:  # a list of commands to run in the container
    - go build -o nut
    - echo Done
  run: # call this macro with "nut run"
    usage: run the project in the container
    actions:
    - ./nut
  test:
    usage: test the project
    actions:
    - go test

container_working_directory: /go/src/project # where macros will be executed
syntax_version: "4" # Nut evolves quickly ; its configuration file syntax as well.
                    # So nut files are versioned to ensure backward compatibility.

```

Here are other instructive examples:
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

### What the Nut???
- build [Nut](https://github.com/matthieudelaro/nut/blob/master/nut.yml) within Nut (never installed Go, and never going to :)
- build [Docker](https://github.com/matthieudelaro/nut/blob/master/examples/docker/nut.yml)
- build and run [Caffe](https://github.com/matthieudelaro/nut/blob/master/examples/caffe/nut.yml) with `nut build`, `nut test`, `nut train-mnist`.
- compile CUDA code on a Mac Book Air, which hasn't got any Nvidia GPU. Just `nut build`
- test code in a whole infrastructure, by defining a macro running *docker-compose* in a container.
- run linux [graphical applications](https://github.com/matthieudelaro/nut/blob/master/examples/geary/nut.yml) on your Mac after installing [XQuartz](http://www.xquartz.org/):

![Linux application on your Mac](https://camo.githubusercontent.com/b32c086f7da89f3365062f9a6a49b7f64377cb35/687474703a2f2f692e696d6775722e636f6d2f4b6650676d72322e676966)

### Milestones
- create container only once, and store its ID in .nut file
- add support for GPU (--device)
- add support for Windows
- add support for *Docker Toolbox* on Mac
- add support for port bindings
- plugin for Sublime Text, to call `nut run`, `nut build`, and `nut test` from the editor.

### Stay Tune
Wanna receive updates? Or share your thoughts? You can post an issue or follow me on [Twitter](https://twitter.com/matthieudelaro).

### Authors and Contributors
@matthieudelaro and @gdevillele