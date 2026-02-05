package com.opstree.microservice.salary.service;

import com.opstree.microservice.salary.entity.SalaryDef;
import java.util.List;
import org.springframework.data.elasticsearch.core.ElasticsearchRestTemplate;
import org.springframework.stereotype.Service;

@Service
public class ElasticsearchRestTemplateServiceImpl implements SalaryService {

    private final ElasticsearchRestTemplate elasticsearchRestTemplate;

    // Manual constructor to bypass Lombok compilation issues
    public ElasticsearchRestTemplateServiceImpl(ElasticsearchRestTemplate elasticsearchRestTemplate) {
        this.elasticsearchRestTemplate = elasticsearchRestTemplate;
    }

    @Override
    public List<SalaryDef> getSalary() {
        // Your existing logic here, for example:
        return null; 
    }
}
