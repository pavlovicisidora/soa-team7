package com.example.tour.service;

import com.example.tour.model.KeyPoint;
import com.example.tour.repository.KeyPointRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Optional;

@Service
public class KeyPointService {
    @Autowired
    KeyPointRepository keyPointRepository;
    public List<KeyPoint> createKeyPoints(List<KeyPoint> keyPoints){
        return keyPointRepository.saveAll(keyPoints);
    }
    public KeyPoint createKeyPoint(KeyPoint keyPoint){return keyPointRepository.save(keyPoint);}
    public List<KeyPoint> findByTourId(Integer id){return keyPointRepository.findByTourId(id);}
    public KeyPoint updateKeyPoint(Integer id, KeyPoint keyPointDetails) {
        Optional<KeyPoint> optionalKeyPoint = keyPointRepository.findById(id);

        if (optionalKeyPoint.isPresent()) {
            KeyPoint existingKeyPoint = optionalKeyPoint.get();

            existingKeyPoint.setName(keyPointDetails.getName());
            existingKeyPoint.setDescription(keyPointDetails.getDescription());
            existingKeyPoint.setLatitude(keyPointDetails.getLatitude());
            existingKeyPoint.setLongitude(keyPointDetails.getLongitude());
            existingKeyPoint.setImageUrl(keyPointDetails.getImageUrl());

            return keyPointRepository.save(existingKeyPoint);
        } else {
            return null;
        }
    }

    public boolean deleteKeyPoint(Integer id){
        if (keyPointRepository.existsById(id)) {
            keyPointRepository.deleteById(id);
            return true;
        }else{
            return false;
        }
    }
}
