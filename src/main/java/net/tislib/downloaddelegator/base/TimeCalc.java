package net.tislib.downloaddelegator.base;

import lombok.extern.log4j.Log4j2;

import java.time.Duration;
import java.time.Instant;
import java.util.concurrent.atomic.AtomicLong;

@Log4j2
public class TimeCalc {

    private Instant start = Instant.now();
    private Instant current = Instant.now();

    public TimeCalc() {
        reset();
    }

    private AtomicLong lastCount = new AtomicLong();
    private long counter = 0;

    public synchronized void reset() {
        start = Instant.now();
        current = Instant.now();
    }

    public synchronized void printSpeedStep(int exceedMillis, long counter) {
        if (Instant.now().toEpochMilli() - current.toEpochMilli() > exceedMillis) {
            printSpeedStep(counter);
        }
    }

    public synchronized void printSpeedStep(int exceedMillis) {
        counter++;
        printSpeedStep(exceedMillis, counter);
    }

    public synchronized void printSpeedStep() {
        printSpeedStep(3000);
    }

    private void printSpeedStep(long counter) {
        long diff = counter - this.lastCount.get();
        long diffMillis = Instant.now().toEpochMilli() - current.toEpochMilli();
        long diffTotalMillis = Instant.now().toEpochMilli() - start.toEpochMilli();
        float lastSpeed = ((float) diff * 1000 / (float) diffMillis);
        float speed = ((float) counter * 1000 / (float) diffTotalMillis);
        log.info(String.format("%.2f ops, %.2f aops %d %n", lastSpeed, speed, counter));
        current = Instant.now();
        this.lastCount.set(counter);
    }
}
