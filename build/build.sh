#!/bin/zsh

if [[ $(pwd | sed 's/\// /g' | awk '{print $NF}') != "app" ]]
then
    echo "wrong directory"
    exit 1
fi


os=(
    linux
    darwin
    windows
    android
    dragonfly
    freebsd
    netbsd
    openbsd
    plan9
    solaris
)

linux=(
    amd64
    386
    arm
    arm64
    ppc64
    ppc64le
    mips
    mipsle
    mips64
    mips64le
)

darwin=(
    amd64
    arm64
    # 386 # unsupported
    # arm # unsupported 
)

windows=( amd64 386 )

android=( arm64 )

dragonfly=( amd64 )

freebsd=( amd64 386 arm )

netbsd=( amd64 386 arm )

openbsd=( amd64 386 arm )

plan9=( amd64 386 )

solaris=( amd64 )


for (( a=1 ; a<=${#os[@]} ; a++ ))
do
    for (( b=1 ; b<=${(P)#os[a][@]} ; b++ ))
    do
        echo -n "Building for ${os[a]}_${(P)os[a][b]} -- "
        mkdir -p ../build/${os[a]}/
        env GOOS=${os[a]} GOARCH=${(P)os[a][b]} go build -o ../build/${os[a]}/hashclock_${os[a]}_${(P)os[a][b]} .
        if [[ -f ../build/${os[a]}/hashclock_${os[a]}_${(P)os[a][b]} ]]
        then 
            echo "OK"

            if [[ ${os[a]} == "windows" ]]
            then
                mv ../build/${os[a]}/hashclock_${os[a]}_${(P)os[a][b]} ../build/${os[a]}/hashclock_${os[a]}_${(P)os[a][b]}.exe
            fi

        else 
            echo "Failed"
        fi

    done
done

