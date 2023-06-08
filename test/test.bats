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

@test "Run framed" {
    run framed --help
    assert_success
    assert_output --partial 'FRAMED (Files and Directories Reusability, Architecture, and Management)'
}

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

@test "Command when file not exist" {
    for command in "${commands_with_template[@]}"
    do
        run framed $command
        assert_failure
        assert_output --partial '☠️ Can not read file'
    done
}

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

@test "Validate spot missing file" {
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

@test "Validate spot missing dir" {
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

@test "Validate spot missing file and dir" {
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

# Here should be tests for allowChildren, allowedPatterns, forbiddenPatterns, maxDepth, minCount, maxCount

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