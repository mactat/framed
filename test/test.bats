setup() {
    # Load
    load 'test_helper/bats-support/load'
    load 'test_helper/bats-assert/load'
    load 'test_helper/bats-file/load'

    # Vars
    TMP_DIR="/tmp/framed-test"
    VALID_URL="https://raw.githubusercontent.com/mactat/framed/master/examples/python.yaml"
    INVALID_URL="https://raw.githubusercontent.com/mactat/framed/master/examples/invalid.yaml"
    DIR="$( cd "$( dirname "$BATS_TEST_FILENAME" )" >/dev/null 2>&1 && pwd )"

    # Commands
    mkdir -p $TMP_DIR
    cd $TMP_DIR
    all_commands=(import visualize create capture verify)
    commands_with_template=(visualize create verify)
}

teardown() {
    sudo rm -R $TMP_DIR
}

### Basic

@test "Run framed" {
    run framed --help
    assert_success
    assert_output --partial 'FRAMED (Files and Directories Reusability, Architecture, and Management)'
}


### Import

@test "Import from valid example" {

    run framed import --example python
    assert_success
    assert_output --partial '✅ Imported successfully'
    assert_file_exists "$TMP_DIR/framed.yaml"

    # assert that content is correct
    run cat $TMP_DIR/framed.yaml
    assert_output --partial 'name: python_example'
}

@test "Import from invalid example" {
    run framed import --example invalid
    assert_failure
    assert_output --partial '☠️ Can not find correct structure'
}

@test "Import from valid url" {
    run framed import --url $VALID_URL
    assert_success
    assert_output --partial '✅ Imported successfully'
    assert_file_exists "$TMP_DIR/framed.yaml"

    # assert that content is correct
    run cat $TMP_DIR/framed.yaml
    assert_output --partial 'name: python_example'
}

@test "Import from invalid url" {
    run framed import --url $INVALID_URL
    assert_failure
    assert_output --partial '☠️ Can not find correct structure'
}

@test "Import with specified output" {
    run framed import --example python --output $TMP_DIR/test.yaml
    assert_success
    assert_output --partial '✅ Imported successfully'
    assert_file_exists "$TMP_DIR/test.yaml"

    # assert that content is correct
    run cat $TMP_DIR/test.yaml
    assert_output --partial 'name: python_example'
}


### Visualize

@test "Visualize show correct structure" {
    # setup
    cp $DIR/test.yaml $TMP_DIR/framed.yaml
    assert_file_exists "$(pwd)/framed.yaml"

    # test
    run framed visualize
    assert_success
    assert_output --partial '✅ Read structure'
    assert_output --partial 'test1'
    assert_output --partial 'test2'
    assert_output --partial 'test3'
    assert_output --partial 'test4'
}


### Create

@test "Create creates correct dirs" {
    # setup
    cp $DIR/test.yaml $TMP_DIR/framed.yaml
    assert_file_exists "$(pwd)/framed.yaml"

    # test
    run framed create
    assert_success
    assert_output --partial '✅ Read structure'
    assert_output --partial '📂 Creating directory'
    assert_dir_exists "$(pwd)/test2"
    assert_dir_exists "$(pwd)/test3"
}

@test "Create creates correct files in dirs" {
    # setup
    cp $DIR/test.yaml $TMP_DIR/framed.yaml
    assert_file_exists "$(pwd)/framed.yaml"

    # test
    run framed create --files
    assert_success
    assert_output --partial '✅ Read structure'
    assert_output --partial '📂 Creating directory'
    assert_output --partial '📄 Creating file'
    assert_dir_exists "$(pwd)/test2"
    assert_dir_exists "$(pwd)/test3"
    assert_file_exists "$(pwd)/test1.md"
    assert_file_exists "$(pwd)/test3/test4.md"
}

@test "Validate validates correct structure" {
    # setup
    cp $DIR/test.yaml $TMP_DIR/framed.yaml
    assert_file_exists "$(pwd)/framed.yaml"
    run framed create --files
    assert_success

    # test
    run framed verify
    assert_success
    assert_output --partial '✅ Read structure'
    assert_output --partial '✅ Verified successfully!'
}


### Validate

@test "Verify should spot missing file" {
    # setup
    cp $DIR/test.yaml $TMP_DIR/framed.yaml
    assert_file_exists "$(pwd)/framed.yaml"
    run framed create --files
    assert_success
    rm "$(pwd)/test1.md"

    # test
    run framed verify
    assert_failure
    assert_output --partial '❌ File not found'
    assert_output --partial 'test1.md'
}

@test "Verify should spot missing dir" {
    # setup
    cp $DIR/test.yaml $TMP_DIR/framed.yaml
    assert_file_exists "$(pwd)/framed.yaml"
    run framed create --files
    assert_success
    rm -R "$(pwd)/test2"

    # test
    run framed verify
    assert_failure
    assert_output --partial '❌ Directory not found'
    assert_output --partial 'test2'
}

@test "Verify should spot missing file and dir" {
    # setup
    cp $DIR/test.yaml $TMP_DIR/framed.yaml
    assert_file_exists "$(pwd)/framed.yaml"
    run framed create --files
    assert_success
    rm "$(pwd)/test1.md"
    rm -R "$(pwd)/test2"

    # test
    run framed verify
    assert_failure
    assert_output --partial '❌ File not found'
    assert_output --partial 'test1.md'
    assert_output --partial '❌ Directory not found'
    assert_output --partial 'test2'
}

@test "Verify should spot wrong pattern if allowedPatterns is set" {
    # setup
    cp $DIR/test.yaml $TMP_DIR/framed.yaml
    assert_file_exists "$(pwd)/framed.yaml"
    run framed create --files
    assert_success

    yq -i '.structure.dirs[0].allowedPatterns[0] = "md"' "$(pwd)/framed.yaml"
    touch "$(pwd)/test2/test3.txt"

    # test
    run framed verify
    assert_failure
    assert_output --partial '❌ Not all files match required pattern'
    assert_output --partial 'test2'
    assert_output --partial 'md'
}

@test "Verify should spot wrong pattern if forbiddenPatterns is set" {
    # setup
    cp $DIR/test.yaml $TMP_DIR/framed.yaml
    assert_file_exists "$(pwd)/framed.yaml"
    run framed create --files
    assert_success

    yq -i '.structure.dirs[0].forbiddenPatterns[0] = "md"' "$(pwd)/framed.yaml"
    touch "$(pwd)/test2/test3.md"

    # test
    run framed verify
    assert_failure
    assert_output --partial '❌ Forbidden pattern (md) matched'
    assert_output --partial 'test2'
}

@test "Verify should spot if there is less files than minCount" {
    # setup
    cp $DIR/test.yaml $TMP_DIR/framed.yaml
    assert_file_exists "$(pwd)/framed.yaml"
    run framed create --files
    assert_success

    yq -i '.structure.dirs[0].minCount = 2' "$(pwd)/framed.yaml"

    # test
    run framed verify
    assert_failure
    assert_output --partial '❌ Min count (2) not met'
    assert_output --partial 'test2'
}

@test "Verify should spot if there is more files than maxCount" {
    # setup
    cp $DIR/test.yaml $TMP_DIR/framed.yaml
    assert_file_exists "$(pwd)/framed.yaml"
    run framed create --files
    assert_success

    yq -i '.structure.dirs[0].maxCount = 1' "$(pwd)/framed.yaml"
    touch "$(pwd)/test2/test3.md"
    touch "$(pwd)/test2/test4.md"

    # test
    run framed verify
    assert_failure
    assert_output --partial '❌ Max count (1) exceeded'
    assert_output --partial 'test2'
}

@test "Verify should spot when maxDepth is exceeded" {
    # setup
    cp $DIR/test.yaml $TMP_DIR/framed.yaml
    assert_file_exists "$(pwd)/framed.yaml"
    run framed create --files
    assert_success

    yq -i '.structure.dirs[0].maxDepth = 1' "$(pwd)/framed.yaml"
    mkdir "$(pwd)/test2/test3"

    # test
    run framed verify
    assert_failure
    assert_output --partial '❌ Max depth exceeded (1)'
    assert_output --partial 'test2'
}

@test "Verify should spot when there are subdirs with children not allowed" {
    # setup
    cp $DIR/test.yaml $TMP_DIR/framed.yaml
    assert_file_exists "$(pwd)/framed.yaml"
    run framed create --files
    assert_success

    yq -i '.structure.dirs[0].allowChildren = false' "$(pwd)/framed.yaml"
    mkdir "$(pwd)/test2/test3"

    # test
    run framed verify
    assert_failure
    assert_output --partial '❌ Children not allowed'
    assert_output --partial 'test2'
}


### Capture

@test "Capture should capture correct structure" {
    # setup
    mkdir "$(pwd)/test1"
    touch "$(pwd)/test1/test2.md"

    # test
    run framed capture
    assert_success
    assert_output --partial '📝 Name:                             default'
    assert_output --partial '📂 Directories:                      1'
    assert_output --partial '📄 Files:                            1'
    assert_output --partial '🔁 Patterns:                         1'
    assert_output --partial '✅ Exported to file'

    # verify
    assert_file_exists "$(pwd)/framed.yaml"
    run cat "$(pwd)/framed.yaml"
    assert_success
    assert_output --partial 'name: default'
    assert_output --partial 'test1'
    assert_output --partial 'test2.md'
    assert_output --partial '.md'
}


### All

@test "All commands fails when file not exist" {
    for command in "${commands_with_template[@]}"
    do
        run framed $command
        assert_failure
        assert_output --partial '☠️ Can not read file'
    done
}