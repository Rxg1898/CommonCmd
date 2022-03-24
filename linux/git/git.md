# git



## config配置

### 查看所有配置

```
git config -l
```

### 仓库配置 

```
// 查看
git config --local -l

// 编辑
git config --local -e
```

### 用户配置

```
// 查看
git config --global -l

// 编辑
git config --global -e

git config --global user.email "you@example.com"
git config --global user.name "Your Name"
```

### 系统配置

```
// 查看
git config --system -l

// 编辑
git config --system -e
```

