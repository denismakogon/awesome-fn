package com.fn.sentiments;

import org.json.JSONObject;


public class HelloFunction {

    private SentimentAnalyzer sentimentAnalyzer = new SentimentAnalyzer();

    public String handleRequest(String input) {
        SentimentAnalyzer sentimentAnalyzer = new SentimentAnalyzer();
        sentimentAnalyzer.initialize();
        SentimentResult sentimentResult = sentimentAnalyzer.analyze(input);
        SentimentClassification sm = sentimentResult.getSentimentClass();

        System.err.println("Sentiment Score: " + sentimentResult.getSentimentScore());
        System.err.println("Sentiment Type: " + sentimentResult.getSentimentType());
        System.err.println("Very positive: " + sentimentResult.getSentimentClass().getVeryPositive()+"%");
        System.err.println("Positive: " + sentimentResult.getSentimentClass().getPositive()+"%");
        System.err.println("Neutral: " + sentimentResult.getSentimentClass().getNeutral()+"%");
        System.err.println("Negative: " + sentimentResult.getSentimentClass().getNegative()+"%");
        System.err.println("Very negative: " + sentimentResult.getSentimentClass().getVeryNegative()+"%");


        JSONObject obj = new JSONObject();
        obj.put("sentiment_score", sentimentResult.getSentimentScore());
        obj.put("sentiment_type", sentimentResult.getSentimentType());
        obj.put("very_positive_probability", sm.getVeryPositive());
        obj.put("positive_probability", sm.getPositive());
        obj.put("neutral_probability", sm.getNeutral());
        obj.put("negative", sm.getNegative());
        obj.put("very_negative_probability", sm.getVeryNegative());



        return obj.toString();
    }
}
