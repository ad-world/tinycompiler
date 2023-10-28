go build -o tinycompiler
CC="gcc"

function comp {
    BN=$(basename -s .teeny $1)
    TTOUTPUT=sudo ./tinycompiler $1 2>&1
    if [ $? -ne 0 ]; then
        echo "${TTOUTPUT}"
    else
        CCOUTPUT=$(${CC} -o "exec/${BN}" "results/${BN}".c)
        if [ $? -ne 0 ]; then
            echo "${CCOUTPUT}"
        else
            echo "${TTOUTPUT}"
        fi
    fi
}

if [ $# -eq 0 ]; then
    echo "Did not provide an input teeny file."
else
    comp $1
fi