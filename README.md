# (Practice) Toy Language Interpreter - Go

<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [(Practice) Toy Language Interpreter - Go](#practice-toy-language-interpreter-go)
  - [使用方法](#使用方法)
    - [编译](#编译)
    - [进入 REPL 模式（交互模式）](#进入-repl-模式交互模式)
    - [执行指定的脚本](#执行指定的脚本)
  - [代码示例](#代码示例)
    - [右折叠](#右折叠)

<!-- /code_chunk_output -->

练习单纯使用 Go lang 编写简单的 _玩具语言_ 解析器。

> 注：本项目是阅读和学习《Writing An Interpreter In Go》时的随手练习，并无实际用途。程序的原理、讲解和代码的原始出处请移步 https://interpreterbook.com/

## 使用方法

### 编译

`$ go build -o toy`

### 进入 REPL 模式（交互模式）

`$ ./toy`

或者

`$ go run .`

### 执行指定的脚本

`$ ./toy path_to_script_file`

或者

`$ go run . path_to_script_file`

示例：

`$ ./toy examples/01-sum.toy`

## 代码示例

### 右折叠

```js
// 定义（从左往右）折叠函数
//
// * list 是一个数组，比如 [1,2,3]
// * initial 是初始值
// * func 是一个函数，签名为
//   (accumulator, element) -> result

let fold = fn(list, initial, func) {
    let iter = fn(list, accumulator) {
        if (len(list) == 0) {
            accumulator
        } else {
            iter(rest(list), func(accumulator, first(list)));
        }
    };
    iter(list, initial);
};

// 使用折叠函数实现对数组元素求和
let sum = fn(list) {
    fold(
        list,
        0,
        fn(accumulator, element) {
            accumulator + element
        }
    );
};

let n = sum([1, 2, 3, 4, 5]);
puts(n); // 输出 15
```