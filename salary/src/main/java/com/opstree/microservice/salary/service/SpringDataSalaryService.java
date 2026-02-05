package com.opstree.microservice.salary.service;

import com.opstree.microservice.salary.entity.SalaryDef;
import java.util.List;
import org.springframework.data.elasticsearch.core.ElasticsearchRestTemplate;
import org.springframework.stereotype.Service;

@Service
public class ElasticsearchRestTemplateServiceImpl implements SpringDataSalaryService {

    private final ElasticsearchRestTemplate elasticsearchRestTemplate;

    public ElasticsearchRestTemplateServiceImpl(ElasticsearchRestTemplate elasticsearchRestTemplate) {
        this.elasticsearchRestTemplate = elasticsearchRestTemplate;
    }

    @Override
    public List<SalaryDef> getSalary() {
        // Logic for Elasticsearch template
        return null; 
    }
}
