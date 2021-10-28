# l2cache

golang内存+磁盘缓存，扩展bytes.Buffer的能力，当数据量超过某个量级时降级缓存到磁盘上

## 使用方式

    cache, err := cache.New(1024, os.TempDir())
    if err != nil {
        logging.Error("can not create cache: %v", err)
        return
    }
    _, err = cache.Write([]byte("hello world"))
    if err != nil {
        logging.Error("can not write data: %v", err)
        return
    }