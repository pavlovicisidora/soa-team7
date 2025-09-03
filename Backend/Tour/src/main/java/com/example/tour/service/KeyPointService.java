package com.example.tour.service;

import com.example.tour.model.KeyPoint;
import com.example.tour.repository.KeyPointRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class KeyPointService {
    @Autowired
    KeyPointRepository keyPointRepository;
    public List<KeyPoint> createKeyPoint(List<KeyPoint> keyPoints){
        return keyPointRepository.saveAll(keyPoints);
    }
}
