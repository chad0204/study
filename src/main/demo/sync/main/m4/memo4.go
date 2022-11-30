package main

// map的value（res）加一个chan作为通知标志

// 加锁 访问与设置标志, 设置标志很快, 所以锁粒度很小

// 相同key的获取, 只有一个g负责写入, 并在写入完成后通知(close chan), 其他g等待写入完成 (<- chan)。
