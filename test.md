# Test/Example

    git init

    echo 1 > 1.txt
    git stage 1.txt
    git commit -m "commit to master #1"

    echo 2 > 2.txt
    git stage 2.txt
    git commit -m "commit to master #2"

    git checkout -b foo
    bopgit track master

    echo foo > foo.txt
    git stage foo.txt
    git commit -m "commit to foo"

    git log --oneline --graph --all

    # * 053a3ed (HEAD -> foo) commit to foo
    # * c2f0beb (master, refs/bopgit/latest-base-commit/foo) commit to master #2
    # * c772312 commit to master #1

    git checkout master

    echo "2 modified" > 2.txt
    git stage 2.txt
    git commit --amend -a -m "commit to master #2 (modified)"

    git log --oneline --graph --all

    # * 303d6a5 (HEAD -> master) commit to master #2 (modified)
    # | * 053a3ed (foo) commit to foo
    # | * c2f0beb (refs/bopgit/latest-base-commit/foo) commit to master #2
    # |/
    # * c772312 commit to master #1

    git checkout foo
    bopgit update
    git log --oneline --graph --all

    # * edb182f (HEAD -> foo) commit to foo
    # * 303d6a5 (master, refs/bopgit/latest-base-commit/foo) commit to master #2 (modified)
    # * c772312 commit to master #1
