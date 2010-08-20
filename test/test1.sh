#!/bin/sh
../json2csv -k a -i test1.json -o test1.output
diff test1.output test1.expected
if [ $? != 0 ]; then
    echo "output not as expected"
    exit 1;
fi
echo "test1 passes";
rm test1.output;