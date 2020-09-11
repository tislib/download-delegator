package net.tislib.downloaddelegator.test;

import org.apache.http.client.config.RequestConfig;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;

import java.net.*;
import java.util.Enumeration;

public class Test2 {

    public static void main(String[] args) throws Exception {
        NetworkInterface nif = NetworkInterface.getByName("enp0s31f6");
        System.out.println("Starting to using the interface: " + nif.getName());
        Enumeration<InetAddress> nifAddresses = nif.getInetAddresses();

        InetAddress nifAddress = nifAddresses.nextElement();


        RequestConfig config = RequestConfig.custom()
                .setLocalAddress(InetAddress.getByName("172.20.11.46")).build();

        HttpGet httpGet = new HttpGet("http://tisserv.net");
        httpGet.setConfig(config);
        CloseableHttpClient httpClient = HttpClients.createDefault();
        try {
            CloseableHttpResponse response = httpClient.execute(httpGet);
            try {
                //logic goes here
            } finally {
                response.close();
            }
        } finally {
            httpClient.close();
        }
    }

}
