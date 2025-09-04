package com.example.tour.controller;

import com.example.tour.model.KeyPoint;
import com.example.tour.service.KeyPointService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequestMapping(value = "/keypoint",produces = MediaType.APPLICATION_JSON_VALUE)
public class KeyPointController {
    @Autowired
    KeyPointService keyPointService;
    @PostMapping("/create")
    public ResponseEntity<List<KeyPoint>> create(@RequestBody List<KeyPoint> keyPoints){
        List<KeyPoint> keyPointsData = keyPointService.createKeyPoints(keyPoints);
        return new ResponseEntity<>(keyPointsData, HttpStatus.OK);
    }
}
