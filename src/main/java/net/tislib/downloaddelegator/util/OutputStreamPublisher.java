package net.tislib.downloaddelegator.util;

import lombok.Getter;
import lombok.SneakyThrows;
import lombok.extern.log4j.Log4j2;
import org.springframework.core.io.buffer.DataBuffer;
import org.springframework.core.io.buffer.DataBufferFactory;
import org.springframework.core.io.buffer.DefaultDataBufferFactory;
import reactor.core.Disposable;
import reactor.core.publisher.Flux;
import reactor.core.publisher.UnicastProcessor;

import java.io.IOException;
import java.io.OutputStream;
import java.time.Duration;
import java.util.function.Consumer;
import java.util.function.Function;

@Log4j2
public class OutputStreamPublisher<T extends OutputStream> {

    @Getter
    private final UnicastProcessor<DataBuffer> publisher = UnicastProcessor.create();

    private final UnicastProcessor<Byte> bufferPublisher = UnicastProcessor.create();
    DataBufferFactory dataBufferFactory = new DefaultDataBufferFactory();

    @Getter
    private final T stream;
    private Disposable subs;

    public OutputStreamPublisher(Function<OutputStream, T> outputStreamFunction) {

        this.stream = outputStreamFunction.apply(new OutputStream() {
            @Override
            public void write(int i) throws IOException {
                bufferPublisher.onNext((byte) i);
            }

            @Override
            public void close() throws IOException {
                super.close();

                bufferPublisher.onComplete();
            }
        });
    }


    public <R> Flux<DataBuffer> mapper(Flux<R> res, Consumer<R> consumer) {

        Flux<R> flux = res.map(item -> {
            consumer.accept(item);

            return item;
        });


        return Flux.just(new Object())
                .map(item -> {
                    subs = flux
                            .doOnComplete(this::closeResponse)
                            .doOnCancel(this::closeResponse)
                            .doOnError(error -> {
                                log.error("ERROR ON STREAM ", error);
                                closeResponse();
                            })
                            .subscribe();

                    return item;
                })
                .flatMap(item -> bufferPublisher)
                .bufferTimeout(4096, Duration.ofMillis(100))
                .map(item -> {
                    Byte[] resp = item.toArray(new Byte[0]);
                    byte[] resp2 = new byte[resp.length];

                    for (int i = 0; i < resp.length; i++) {
                        resp2[i] = resp[i];
                    }

                    return dataBufferFactory.wrap(resp2);
                })
                .doOnCancel(() -> {
                    subs.dispose();
                });
    }

    @SneakyThrows
    private void closeResponse() {
        stream.close();
    }
}
