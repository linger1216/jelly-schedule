
# gradleSrc=./

# function updateVersionHomePatch(){
#     touch ./version.meta
#     grep -o "def versionHomePatch = \([0-9]*\)" $gradleSrc > ./version.meta
#     version1=`sed 's/.*def versionHomePatch = \([0-9]*\)/\1/' ./version.meta`
#     echo "当前版本：$version1";
#     add=1
#     count=$[$version1+add]
#     echo "打包版本：$count"
#     new="def versionHomePatch = $count"
#     sed -i "" "s/def versionHomePatch = $version1/def versionHomePatch = $count/" 
#     cd $gradleSrc
#     rm -f ./version.meta
# }

# updateVersionHomePatch

# 通过grep获取版本号的位置，存到临时文件tmpshelldata.txt
# 使用sed提取tmpshelldata.txt中的具体版本号，并且用命令替换赋值给一个变量version1
# 通过 version1 加 1，计算新的版本号
# 使用sed命令，用新版本号 代替 老版本号
# 删除临时文件

#rm -fr build
#mkdir -p build/bin

go build -o build/bin/api cmd/api/api.go
go build -o build/bin/executor cmd/executor/executor.go
go build -o build/bin/echo-job example/echo-job/main.go
go build -o build/bin/shell-job example/shell-job/main.go
go build -o build/bin/http-job example/http-job/main.go

# locker
#go build -o build/bin/locker1 cmd/locker/*.go
#go build -o build/bin/locker2 cmd/locker/*.go
#go build -o build/bin/locker3 cmd/locker/*.go
#go build -o build/bin/locker4 cmd/locker/*.go