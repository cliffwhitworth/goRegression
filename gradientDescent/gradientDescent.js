const data = require("../data/advertising.json");
var _ = require('./lodash.js');

const x = [];
const y = [];
// const selectedFeatures = ['TV', 'Radio', 'Newspaper'];
const selectedFeatures = ['TV'];
const selectedLabels = ['Sales'];

Object.keys(data).forEach(function(key, keyIndex) {
    var temp = [];
    Object.keys(data[key])
      .filter(key => selectedFeatures.includes(key))
      .reduce((obj, kee) => {
        temp.push(data[key][kee]);
    }, {});
    x.push(temp);
});

Object.keys(data).forEach(function(key, keyIndex) {
    var temp = [];
    Object.keys(data[key])
      .filter(key => selectedLabels.includes(key))
      .reduce((obj, kee) => {
        temp.push(data[key][kee]);
    }, {});
    y.push(temp);
});

class GradientDescentJS {
    constructor(features, labels, options) {
        this.features = features;
        this.labels = labels;
        this.mseHistory = [];
        this.N = features.length;

        this.getStats();
    
        // assign default values to options
        this.options = Object.assign(
          { learningRate: 0.1, iterations: 10 },
          options
        );
        
        // standardize
        // uncomment yHat and mSlope
        this.features = this.features.map((val, i) => {
            return (val - this.xMean) / Math.pow(this.uxVariance, 0.5);
        });
    
        this.b = 0;
        this.m = 0;
    }

    gradientDescent() {
        const yHat = this.features.map(row => {
            // return this.m * row[0] + this.b
            return this.m * row + this.b
        });

        this.mse = _.sum(yHat.map((val, i) => {
            return Math.pow((this.labels[i][0] - val), 2) / (this.labels.length);
        }));

        this.mseHistory.unshift(this.mse);
        
        const bSlope = (_.sum(yHat.map((guess, i) => {
            return guess - this.labels[i][0]
        })) * 2) / this.features.length;

        const mSlope = (_.sum(yHat.map((guess, i) => {
            // return -1 * this.features[i][0] * (this.labels[i][0] - guess)
            return -1 * this.features[i] * (this.labels[i][0] - guess);
        })) * 2) / this.features.length;

        this.m = this.m - mSlope * this.options.learningRate;
        this.b = this.b - bSlope * this.options.learningRate;
    }

    train() {
        for( let i = 0; i < this.options.iterations; i++ ) {
            this.gradientDescent();
            this.updateLearningRate();
        }        
    }
    
    updateLearningRate() {
        if (this.mseHistory.length < 2) {
            return;
        }
    
        if (this.mseHistory[0] > this.mseHistory[1]) {
            this.options.learningRate /= 2;
        } else {
            this.options.learningRate *= 1.05;
        }
    }

    getStats() {
        // stats for x
        this.xMean = _.sum(this.features.map((val, i) => {
            return val / this.features.length;
        }));

        this.xVariance = _.sum(this.features.map((val, i) => {
            return Math.pow((val - this.xMean), 2) / this.features.length;
        }));

        this.uxVariance = _.sum(this.features.map((val, i) => {
            return Math.pow((val - this.xMean), 2) / (this.features.length - 1);
        }));

        // stats for y
        this.yMean = _.sum(this.labels.map((val, i) => {
            return val / this.labels.length;
        }));

        this.yVariance = _.sum(this.labels.map((val, i) => {
            return Math.pow((val - this.yMean), 2) / this.labels.length;
        }));

        this.uyVariance = _.sum(this.labels.map((val, i) => {
            return Math.pow((val - this.yMean), 2) / (this.labels.length - 1);
        }));
    }
}

regression = new GradientDescentJS(x, y, { learningRate: 0.1, iterations: 100 })
console.log("N: ", regression.N)
console.log("xMean: ", regression.xMean);
console.log("yMean: ", regression.yMean);
console.log("xVariance: ", regression.xVariance);
console.log("unbiased xVariance: ", regression.uxVariance);
console.log("yVariance: ", regression.yVariance);
console.log("unbiased yVariance: ", regression.uyVariance);
regression.train()
console.log("slope: ", regression.m)
console.log("bias: ", regression.b)