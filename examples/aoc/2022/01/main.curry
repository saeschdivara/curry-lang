package main

import (
    "internal/os"
);

// resolve example
fn resolveAssignment(filePath) {
    let file = os.open(filePath);
    let fileContent = file.readAll();
    
    let maxVal = fileContent.split("\n\n")
                            .map(fn(p) { 
                                p.split("\n")
                                 .map(fn (l) { l.toInt() })
                                 .sum()
                            })
                            .max();
    
    os.Printf("Max val: %d", maxVal);
}

resolveAssignment("data.example");