---
title: Go:迭代器
date: 2026-06-11 15:36:21
tags: Go
categories: GO
---

​	由浅入深理解迭代器。

<!--more-->

### 什么是迭代器

​	简单说，迭代器是一个函数，可以直接用 **for range** 遍历。函数的签名应该为以下之一：

```go
	func(func() bool)			//无值迭代器
	func(func(V) bool)			//单值迭代器
	func(func(K, V) bool)		//两值迭代器
```

### 最简单无值迭代器

​	定义：

```go
    // 无返回值迭代器，仅迭代次数
    func iterator0(yield func() bool) {
        for i := 0; i < 10; i++ {
            if !yield() {
                return
            }
        }
    }
```

​	使用：

```go
        m := 0
        for range iterator0 {
            fmt.Println(m)//打印0到9
            m++
        }

```

### 控制迭代次数

​	为了能控制迭代次数，我们需要一个 **生成迭代器** 的函数：

```go
    // 获得一个 无返回值迭代器
    func iterator1(n int) func(func() bool) {
        return func(yield func() bool) {
            for i := 0; i < n; i++ {
                if !yield() {
                    return
                }
            }
        }
    }
```

​	使用时，输入参数控制次数：

```go
	n := 0
	for range iterator1(7) {
		fmt.Println(n)//打印0到6
		n++
	}
```

​	对比一下两者的使用方式：

```go
	for range iterator0
	for range iterator1(7)
```

### 单值迭代器

​	无值迭代器只是展示写法，实际中用不到。实际使用迭代器的场景需要迭代器返回数据，比如：

```go
    // 输出一个值的迭代器
    func iterator2(yield func(v int) bool) {
        for i := 0; i < 7; i++ {
            if !yield(i) {
                return
            }
        }
    }
	
	//打印0到6
	for i := range iterator2 {
		fmt.Println(i)
	}
```

​	更灵活：

```go
    // 获得一个 输出一个值的迭代器
    func iterator3(n int) func(func(int) bool) {
        return func(yield func(v int) bool) {
            for i := 0; i < n; i++ {
                if !yield(i) {
                    return
                }
            }
        }
    }
	
	//打印0到3
    for i := range iterator3(4) {
            fmt.Println(i)
        }
```

### 迭代切片/map

​	**iterator3** 只控制了迭代值的上限，更正常的情况是遍历一个切片：

```go
    func iterator4(slice []int) func(func(int) bool) {
        return func(yield func(int) bool) {
            for i := 0; i < len(slice); i++ {
                if !yield(slice[i]) {
                    return
                }
            }
        }
    }
	
	//打印切片
	arr := []int{2, 8, 6, 3}
	for i := range iterator4(arr) {
		fmt.Println(i)
	}
```

### 双值迭代器

```go
    func iterator5(slice []int) func(func(i, v int) bool) {
        return func(yield func(i, v int) bool) {
            for i := 0; i < len(slice); i++ {
                if !yield(i, slice[i]) {
                    return
                }
            }
        }
    }
    func iterator6(m map[int]string) func(func(i int, v string) bool) {
        return func(yield func(i int, v string) bool) {
            for k, v := range m {
                if !yield(k, v) {
                    return
                }
            }
        }
    }

	//切片
	for i, v := range iterator5([]int{43, 6, 9, 2}) {
		fmt.Println(i, v)
	}
	//map
	for k, v := range iterator6(map[int]string{1: "one", 2: "two", 3: "three"}) {
		fmt.Println(k, v)
	}
```

### 模板

```go
    type Seq[V any] func(func(V) bool)
    type Seq2[K any, V any] func(func(K, V) bool)

    // 模板 迭代切片
    // ~ 符号表示允许底层类型匹配
    func All[Slice ~[]E, E any](slice Slice) Seq2[int, E] {
        return func(yield func(int, E) bool) {
            for i, v := range slice {
                if !yield(i, v) {
                    return
                }
            }
        }
    }
    // 模板 迭代Map
    func MapAll[Map ~map[K]V, K comparable, V any](m Map) Seq2[K, V] {
        return func(yield func(K, V) bool) {
            for k, v := range m {
                if !yield(k, v) {
                    return
                }
            }
        }
    }
```

​	使用模板：

```go
	// 注意模板可以类型推断
	for i, v := range All([]int{8, 73, 2, 90}) {
		fmt.Println(i, v)
	}
```

​	自定义类型接收器写法：
```go
    type MySlice[T any] []T

    func (s MySlice[T]) All() Seq2[int, T] {
        return func(yield func(int, T) bool) {
            for i, v := range s {
                if !yield(i, v) {
                    return
                }
            }
        }
    }
```

### Pull 迭代器

​	以上列举迭代器皆为 **Push 迭代器**。而 **Pull 迭代器** 由用户主动**拉取**下一个值。

​	标准库 `iter` 包提供了 `Pull` 函数：

```go
	import "iter"
	
	func Pull[V any](seq iter.Seq[V]) (next func() (V, bool), stop func())
	func Pull2[K, V any](seq Seq2[K, V]) (next func() (K, V, bool), stop func())
```

​	如：

```go
	next, stop := iter.Pull2(slices.All([]int{4, 5, 3, 7, 6}))
	i, v, ok := next()
	fmt.Println(i, v, ok)  //4
	i, v, ok = next()
	fmt.Println(i, v, ok)  //5
	stop()
```

### 其他

​	实现两两输出迭代器的两种写法：

```go
    func Pairs1(slice []int) Seq2[int, int] {
        return func(yield func(int, int) bool) {
            for i := 0; i < len(slice) && i+1 < len(slice); i += 2 {
                v1 := slice[i]
                v2 := slice[i+1]
                if !yield(v1, v2) {
                    return
                }
            }
        }
    }
    func Pairs2(seq iter.Seq[int]) Seq2[int, int] {
        return func(yield func(int, int) bool) {
            next, stop := iter.Pull(seq)
            defer stop()
            for {
                v1, ok := next()
                if !ok {
                    return
                }
                v2, ok := next()
                if !ok {
                    return
                }
                if !yield(v1, v2) {
                    return
                }
            }
        }
    }
```

​	随机数迭代器：

```go
    func RandomSeq(aMax, aNum int) iter.Seq[int] {
        return func(yield func(int) bool) {
            for i := 0; i < aNum; i++ {
                v := rand.Intn(aMax)
                if !yield(v) {
                    return
                }
            }
        }
    }
```



### 总结

1. 即便嵌套复杂，核心就是迭代器符合规定的**函数签名**
2. 数据通过 **yield** 传输，且传输的数据没有限制（如上面的两两迭代）



