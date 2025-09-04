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
    public List<KeyPoint> createKeyPoints(List<KeyPoint> keyPoints){
        return keyPointRepository.saveAll(keyPoints);
    }
    public KeyPoint createKeyPoint(KeyPoint keyPoint){return keyPointRepository.save(keyPoint);}
    public List<KeyPoint> findByTourId(Integer id){return keyPointRepository.findByTourId(id);}
}
