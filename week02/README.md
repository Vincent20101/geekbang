~~~
我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
~~~


> dao 返回 sql.ErrNoRows ,是否应该 wrap 这个 error ，分两种情况考虑:

~~~

情况1： dao层的封装，如果是作为第三方库或者公共库提供给其他人使用，或者跨项目使用，这时不应该对该错误进行warp包装。 
因为，如果wrap这个错误，而上层调用时自己也wrap包装了错误的话，就会出现两处堆栈日志冗余打印。如果一定要增加上下文文本信息，但又不附加调用栈，可以使用 WithMessage 方法。 


情况2：如果是本业务应用层使用dao操作时，或者可控的项目范围内时，可以有两种做法：

 ① 如果上层需要关心数据为空的情况，可以使用wrap对错误进行包装，在程序的顶层使用 %+v 打印堆栈的详情记录
 ② 如果上层调用不关心数据是否为空，有没数据都可以的话，可以直接返回空数组
  
 （代码详见：week02/main.go）
~~~