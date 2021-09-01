# prew

prew는 파이썬 프로젝트를 관리하기 위해 만든 툴입니다.

## 프로젝트 관리

prew는 프로젝트 관리를 위한 몇가지 유용한 명령어를 제공합니다.

### init

`prew init <directory>`

해당 경로에 새로운 프로젝트를 생성합니다.  
생성 시 prew는 유저에게 프로젝트 이름과 파이썬 버전을 입력하도록 요청합니다.  
가상 환경의 경우엔 virtualenv(https://github.com/pypa/virtualenv)를 사용합니다.

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

### search

`prew search <package>`

해당 패키지를 pypi에서 검색합니다.  
직후 버전 목록을 출력하며 사용자가 선택한 버전을 설치할 지 물어봅니다.  
설치한다고 응답할 경우 해당 패키지의 해당 버전을 설치합니다.

### install

`prew install <package> [<version>]`

현재 경로의 prew 프로젝트에 해당 패키지를 설치합니다.  
입력받은 패키지의 버전을 설치하고 spec.yaml에 패키지 정보를 저장합니다.  
버전을 입력하지 않으면 최신 버전을 설치합니다.

### list

`prew list`

현재 경로의 prew 프로젝트에 설치되어 있는 패키지를 출력합니다.  
출력된 패키지 중 하나를 선택하면 해당 패키지를 삭제할 지 물어봅니다.  
삭제한다고 응답할 경우 해당 패키지를 삭제합니다.

### remove

`prew remove <package> [<flag>...]`

현재 경로의 prew 프로젝트에 설치된 해당 패키지를 제거합니다.  
만약 존재하지 않을 경우엔 그냥 넘어갑니다.  

`-y` 플래그를 입력 받으면 사용자에게 묻지 않고 제거합니다.  
`-d` 플래그를 입력 받으면 의존성도 삭제합니다.

### tidy

`prew tidy [<flag>]`

현재 경로의 prew 프로젝트의 파이썬 파일과 spec.yaml을 비교하여 사용하지 않는 패키지를 제거합니다.  
이 때 다른 패키지의 의존성으로 설치된 패키지는 제외됩니다.  

`-y` 플래그를 입력 받으면 사용자에게 묻지 않고 제거합니다.

### check

`prew check [<filename>] [<flag>]`

현재 경로의 prew 프로젝트의 특정 파이썬 파일의 타입을 체크합니다.  
filename을 입력하면 해당 파일만 체크하지만 -a 플래그를 입력하면 모든 파이썬 파일을 체크합니다.  
이 기능은 mypy(https://github.com/python/mypy)를 사용합니다.

## third party

using virtualenv(https://github.com/pypa/virtualenv) for virtual environment  
using mypy(https://github.com/python/mypy) to check python file type

## Special thanks

@lemon-mint: github.com/lemon-mint