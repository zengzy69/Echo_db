1.线程模型选择了多线程的锁模型

2.数据一致性算法选择了最终一致性Gossip协议

3.使用了gorm从外部MySQL读取预加载数据

4.使用了B+树的数据结构

5.过期键删除策略选择了定期删除

6.内存淘汰策略选择LFU算法

7.业务基本实现版本号的检查和更新

![image](https://github.com/user-attachments/assets/6ee07a85-91a4-4b2c-9d5f-1b06c812ddec)


![ddff59195e4e3ea18f8dd9c4d3f596e0](https://github.com/user-attachments/assets/6ba8dbdb-38f7-4af6-ad03-895dbc016b28)
