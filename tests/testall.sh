#!/bin/bash

here=$(dirname "$0")

filter="$*"
: ${filter:='*'}

if [[ "$WITH_DEBUGGER" == 1 ]] ; then
    debug=--debug
fi

export options="--read-buffer-size 4096"
export debug

bin=../bin/hl

die() {
  echo "$@"
  exit 2
}

unset COLORTERM

cd $here || die "$0: can't chdir to $here."

num_pass=0
num_fail=0

../scripts/build.sh || exit 1

for r in t*.rules ; do
  name=$(basename $r .rules)

  for input in ${name}_*.input; do
    if [[ "$input" != $filter ]] ; then
        continue
    fi
    echo -n "Test: $input "
    base="$(basename "$input" .input)"
    expect="${base}.expect"
    actual="${base}.actual"
    diff="${base}.diff"
    bash $r <"$input" >"$actual"
    diff --color=never -u "$expect" "$actual" >"$diff"
    rc=$?
    if (( $rc == 0 )) ; then
      echo "pass"
      num_pass=$(( $num_pass + 1 ))
    else
      echo "FAIL"
      num_fail=$(( $num_fail + 1 ))
    fi
  done
done

if (( $num_pass > 0 && $num_fail == 0 )) ; then
    exit 0
else
    exit 1
fi
