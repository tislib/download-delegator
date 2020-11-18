package net.tislib.downloaddelegator.base;

import java.time.Duration;
import java.time.Instant;
import java.util.concurrent.atomic.AtomicLong;

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

    public synchronized void printStep() {
        Instant now = Instant.now();
        long last = Duration.between(current, now).toMillis();
        long full = Duration.between(start, now).toMillis();
        if (last > 10000) {
            last /= 1000;
            System.out.print("Time taken: " + last + " seconds last step; ");
        } else {
            System.out.print("Time taken: " + last + " milliseconds last step; ");
        }

        if (full > 10000) {
            full /= 1000;
            System.out.println(full + " seconds last beginning");
        } else {
            System.out.println(full + " milliseconds last beginning");
        }
        current = now;
    }

    public synchronized void runIfExceed(long millis, Runnable runnable) {
        if (Instant.now().toEpochMilli() - millis > current.toEpochMilli()) {
            runnable.run();
            current = Instant.now();
        }
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
        System.out.printf("%.2f ops, %.2f aops %d %n", lastSpeed, speed, counter);
        current = Instant.now();
        this.lastCount.set(counter);
    }

    public void printSpeedStep(int exceedMillis, long counter, int size) {
        if (Instant.now().toEpochMilli() - current.toEpochMilli() > exceedMillis) {
            printSpeedStep(counter);
            System.out.printf("total: %d; percentage: %.2f%% %n", size, ((float) counter / size) * 100);
        }
    }
}
