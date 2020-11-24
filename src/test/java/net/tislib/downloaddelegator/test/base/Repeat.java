package net.tislib.downloaddelegator.test.base;

import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

@Retention(RetentionPolicy.RUNTIME)
@Target({
        java.lang.annotation.ElementType.METHOD
})
public @interface Repeat {
    public abstract int times();
}
