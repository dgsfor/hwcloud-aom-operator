# hwcloud-aom-operator

分为两个模块：

- aom group controller
- aom policy controller

### group controller
主要是提供deployment的aom组策略，包含三个参数：

- MaxInstances 伸缩pod上限
- MinInstances 伸缩pod下限
- CooldownTime 两个伸缩活动之间的冷却时间

### 注意

使用的时候，请修改一下内容：

```shell
1. 重命名`config/const/prod-git.ini`为`config/const/prod.ini`
2. 修改`config/const/prod.ini`的内容
```