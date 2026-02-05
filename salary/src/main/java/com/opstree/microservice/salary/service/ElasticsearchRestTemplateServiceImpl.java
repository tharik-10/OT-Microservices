package com.opstree.microservice.salary.service;

import com.opstree.microservice.salary.entity.SalaryDef;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;
import lombok.extern.slf4j.Slf4j;
import org.elasticsearch.action.index.IndexRequest;
import org.elasticsearch.client.indices.CreateIndexRequest;
import org.elasticsearch.common.settings.Settings;
import org.elasticsearch.index.query.QueryBuilders;
import org.elasticsearch.search.aggregations.AggregationBuilder;
import org.elasticsearch.search.aggregations.bucket.terms.Terms;
import org.elasticsearch.search.aggregations.bucket.terms.TermsAggregationBuilder;
import org.springframework.data.elasticsearch.core.ElasticsearchRestTemplate;
import org.springframework.data.elasticsearch.core.SearchHit;
import org.springframework.data.elasticsearch.core.SearchHits;
import org.springframework.data.elasticsearch.core.query.Criteria;
import org.springframework.data.elasticsearch.core.query.CriteriaQuery;
import org.springframework.data.elasticsearch.core.query.NativeSearchQueryBuilder;
import org.springframework.data.elasticsearch.core.query.Query;
import org.springframework.stereotype.Service;

@Service
@Slf4j
// Removed @RequiredArgsConstructor because the build environment is failing to process it
public class ElasticsearchRestTemplateServiceImpl implements ElasticsearchRestTemplateService {

    private final ElasticsearchRestTemplate elasticsearchRestTemplate;

    // MANUAL CONSTRUCTOR: This does exactly what @RequiredArgsConstructor used to do.
    public ElasticsearchRestTemplateServiceImpl(ElasticsearchRestTemplate elasticsearchRestTemplate) {
        this.elasticsearchRestTemplate = elasticsearchRestTemplate;
    }

    public List<SalaryDef> getSalaryInfo() {
        Query query = new NativeSearchQueryBuilder()
                .withQuery(QueryBuilders.matchAllQuery())
                .build();
        SearchHits<SalaryDef> searchHits = elasticsearchRestTemplate.search(query, SalaryDef.class);

        return searchHits.get().map(SearchHit::getContent).collect(Collectors.toList());
    }
}
