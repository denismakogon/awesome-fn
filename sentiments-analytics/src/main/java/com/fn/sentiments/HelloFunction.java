package com.fn.sentiments;

import org.json.JSONObject;


public class HelloFunction {

    private SentimentAnalyzer sentimentAnalyzer = new SentimentAnalyzer();

    public String handleRequest(String input) {
        SentimentAnalyzer sentimentAnalyzer = new SentimentAnalyzer();
        sentimentAnalyzer.initialize();
        SentimentResult sentimentResult = sentimentAnalyzer.analyze(input);
        SentimentClassification sm = sentimentResult.getSentimentClass();

        JSONObject obj = new JSONObject();
        obj.put("sentiment_score", sentimentResult.getSentimentScore());
        obj.put("sentiment_type", sentimentResult.getSentimentType());
        obj.put("very_positive_probability", sm.getVeryPositive());
        obj.put("positive_probability", sm.getPositive());
        obj.put("neutral_probability", sm.getNeutral());
        obj.put("negative", sm.getNegative());
        obj.put("very_negative_probability", sm.getVeryNegative());

        System.err.println(obj.toString());

        return obj.toString();
    }
}
