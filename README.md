# prew

prew는 파이썬 프로젝트를 관리하기 위해 만든 툴입니다.

## 프로젝트 관리

prew는 프로젝트 관리를 위한 몇가지 유용한 명령어를 제공합니다.

### init

`prew init <directory>`

해당 경로에 새로운 프로젝트를 생성합니다.  
생성 시 prew는 유저에게 프로젝트 이름과 파이썬 버전을 입력하도록 요청합니다.

### run

`prew run`

현재 경로에 prew 프로젝트가 있다면 해당 프로젝트를 실행합니다.

### restore

`prew restore`

현재 경로에 prew 프로젝트 파일(spec.yaml)이 있다면 해당 프로젝트의 가상 환경을 만들고 필요한 패키지를 설치합니다.

### make

`prew make [<flags> ...]`

현재 경로에 prew 프로젝트가 있다면 입력된 flag에 따른 결과물을 만듭니다.  
1. -d, --dockerfile: 해당 프로젝트의 dockerfile을 생성합니다.

## 패키지 관리

prew는 pip를 그대로 사용하지만, 패키지를 설치하거나 삭제할 때 사용할 수 있는 편리한 명령어를 제공합니다.

### install

`prew install <package>`

현재 경로의 prew 프로젝트에 해당 패키지를 설치합니다.  
이 명령어가 실행되면 곧 유저에게 패키지의 버전을 선택하도록 요청하고 해당 버전을 설치하면서 spec.yaml에 패키지 정보를 저장합니다.

### remove

`prew remove <package>`

현재 경로의 prew 프로젝트에 설치된 해당 패키지를 제거합니다.  
만약 존재하지 않을 경우엔 그냥 넘어갑니다.

### tidy

`prew tidy`

현재 경로의 prew 프로젝트의 파이썬 파일과 spec.yaml을 비교하여 사용하지 않는 패키지를 제거합니다.  
이 때 다른 패키지의 의존성으로 설치된 패키지는 제외됩니다.

## Special thanks

@lemon-mint: github.com/lemon-mint