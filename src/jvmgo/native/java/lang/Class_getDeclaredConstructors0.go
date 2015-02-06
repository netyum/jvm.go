package lang

import (
    //"fmt"
    //"strings"
    //. "jvmgo/any"
    "jvmgo/jvm/rtda"
    rtc "jvmgo/jvm/rtda/class"
)

// private native Constructor<T>[] getDeclaredConstructors0(boolean publicOnly);
// (Z)[Ljava/lang/reflect/Constructor;
func getDeclaredConstructors0(frame *rtda.Frame) {
    stack := frame.OperandStack()
    publicOnly := stack.PopBoolean()
    jClass := stack.PopRef() // this
    goClass := jClass.Extra().(*rtc.Class)
    goConstructors := goClass.GetConstructors(publicOnly)
    constructorCount := int32(len(goConstructors))
    
    constructorClass := goClass.ClassLoader().LoadClass("java/lang/reflect/Constructor")
    constructorArr := constructorClass.NewArray(constructorCount)
    constructorInitMethod := constructorClass.GetConstructor("(Ljava/lang/Class;[Ljava/lang/Class;[Ljava/lang/Class;IILjava/lang/String;[B[B)V")
    
    stack.PushRef(constructorArr)

    if constructorCount > 0 {
        /*
        Constructor(Class<T> declaringClass,
                    Class<?>[] parameterTypes,
                    Class<?>[] checkedExceptions,
                    int modifiers,
                    int slot,
                    String signature,
                    byte[] annotations,
                    byte[] parameterAnnotations)
        }
        */
        jConstructors := constructorArr.Fields().([]*rtc.Obj)
        thread := frame.Thread()
        for i, goConstructor := range goConstructors {
            jConstructor := constructorClass.NewObj()
            jConstructors[i] = jConstructor

            newFrame := thread.NewFrame(constructorInitMethod)
            vars := newFrame.LocalVars()
            vars.SetRef(0, jConstructor) // this
            vars.SetRef(1, jClass) // declaringClass
            vars.SetRef(2, nil) // todo parameterTypes
            vars.SetRef(3, nil) // todo checkedExceptions
            vars.SetInt(4, int32(goConstructor.GetAccessFlags())) // modifiers
            vars.SetInt(5, int32(0)) // todo slot
            vars.SetRef(6, nil) // todo signature
            vars.SetRef(7, nil) // todo annotations
            vars.SetRef(8, nil) // todo parameterAnnotations
            thread.PushFrame(newFrame)
        }

        //panic("getDeclaredConstructors0"+goClass.Name())
    }
}
