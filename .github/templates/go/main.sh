cd "${PROJECT_PATH}"/99-labs/code || exit 1
if [ "${TESTS}" = "all" ]; then
  items=("helloworld" "splitdim", "kvstore")
  for item in "${items[@]}"; do
    pushd . > /dev/null
    cd "${item}" && apply_and_wait && go mod download &&  go test ./... --tags=kubernetes --count 1
    ret=$?
    if [[ $ret -ne 0 ]]; then
      exit $ret
    fi
    popd > /dev/null
  done
else
  for TEST in ${TESTS}; do
    pushd . > /dev/null
    cd ${TEST} || exit 1
    go mod download
    go test ./... --tags=kubernetes --count 1
    ret=$?
    if [[ $ret -ne 0 ]]; then
      exit $ret
    fi
    popd > /dev/null
  done
fi
