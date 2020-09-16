### 总结了一个解决冲突的常规流程：

1.前提条件：不能在 master 分支上修改任何文件。master 分支的变更只能通过 git pull 和 git merge 获得。在 master 分支下面，不能手动修改任何文件。
2.我们自己有一个分支用来修改代码，例如我的分支叫做dev分支。我把代码修改完成了，现在不知道有没有冲突。
3.在 dev 分支里面，执行命令git merge origin/master，把远程的master分支合并到当前dev分支中。如果没有任何报错，那么直接转到第5步。
4.如果有冲突，根据提示，把冲突解决，保存文件。然后执行命令git add xxx把你修改的文件添加到缓存区。然后执行命令git commit -m "xxx"添加 commit 信息。
5.执行如下命令，切换到 master 分支：git checkout master。
6.执行命令git pull确保当前 master 分支是最新代码。
7.把dev分支的代码合并回 master 分支：git merge dev。
8.提交代码：git push。