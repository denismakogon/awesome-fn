package com.fn.sentiments;

import java.util.Properties;

import org.ejml.simple.SimpleMatrix;

import edu.stanford.nlp.ling.CoreAnnotations;
import edu.stanford.nlp.neural.rnn.RNNCoreAnnotations;
import edu.stanford.nlp.pipeline.Annotation;
import edu.stanford.nlp.pipeline.StanfordCoreNLP;
import edu.stanford.nlp.sentiment.SentimentCoreAnnotations;
import edu.stanford.nlp.trees.Tree;
import edu.stanford.nlp.util.CoreMap;


public class SentimentAnalyzer {

    static Properties props;
    static StanfordCoreNLP pipeline;


    public void initialize() {
        props = new Properties();
        props.setProperty("annotators", "tokenize, ssplit, parse, sentiment");
        pipeline = new StanfordCoreNLP(props);
    }

    public SentimentResult analyze(String text) {


        SentimentResult sentimentResult = new SentimentResult();
        SentimentClassification sentimentClass = new SentimentClassification();

        if (text != null && text.length() > 0) {

            Annotation annotation = pipeline.process(text);

            for (CoreMap sentence : annotation.get(CoreAnnotations.SentencesAnnotation.class)) {
                Tree tree = sentence.get(SentimentCoreAnnotations.SentimentAnnotatedTree.class);
                SimpleMatrix sm = RNNCoreAnnotations.getPredictions(tree);
                String sentimentType = sentence.get(SentimentCoreAnnotations.SentimentClass.class);

                sentimentClass.setVeryPositive(sm.get(4));
                sentimentClass.setPositive(sm.get(3));
                sentimentClass.setNeutral(sm.get(2));
                sentimentClass.setNegative(sm.get(1));
                sentimentClass.setVeryNegative(sm.get(0));

                sentimentResult.setSentimentScore(RNNCoreAnnotations.getPredictedClass(tree));
                sentimentResult.setSentimentType(sentimentType);
                sentimentResult.setSentimentClass(sentimentClass);
            }
        }
        return sentimentResult;
    }
}
