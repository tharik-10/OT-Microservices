package com.opstree.microservice.salary.service;

import com.opstree.microservice.salary.entity.SalaryDef;
import com.opstree.microservice.salary.repository.SalaryRepository;
import java.util.List;
import org.springframework.stereotype.Service;

@Service
public class SpringDataSalaryServiceImpl implements SpringDataSalaryService {

    private final SalaryRepository salaryRepository;

    public SpringDataSalaryServiceImpl(SalaryRepository salaryRepository) {
        this.salaryRepository = salaryRepository;
    }

    @Override
    public List<SalaryDef> getSalary() {
        return salaryRepository.findAllSalary();
    }
}
