from __future__ import division
import numpy as np
from decimal import Decimal





class Indicator(object):
    """ Indicator is the base class for technical indicators like moving
    averages. Subclasses should include a calculate(candles) method which
    defines the mathematics of the indicator, and a plot_type attribute
    which indicates whether it should be plotted on the large upper subplot
    ('primary') or the smaller bottom subplot ('secondary') in the GUI.

    Indicators can be used to define Conditions and can be directly compared
    (e.g. EMA10 > EMA21) if containing value arrays of the same length.
    Comparisons can also be made between Indicators and numbers, e.g.
    MACD > 0.  Comparisons yield a boolean array of the same length. Care
    should be taken when making comparisons this way or using the stored
    values, since they represent only the result of the last calculate(candles)
    call, which may become outdated.
    """
    def __init__(self, label=None):
        """ label will be used in the GUI plot legend.
        """
        self.label = label
        self.values = None

    def __lt__(self, other):
        try:
            return np.array(self.values) < np.array(other.values)
        except:
            return np.array(self.values) < other

    def __gt__(self, other):
        try:
            return np.array(self.values) > np.array(other.values)
        except:
            return np.array(self.values) > other

    def __le__(self, other):
        try:
            return np.array(self.values) <= np.array(other.values)
        except:
            return np.array(self.values) <= other

    def __ge__(self, other):
        try:
            return np.array(self.values) >= np.array(other.values)
        except:
            return np.array(self.values) >= other

    def __eq__(self, other):
        try:
            return np.array(self.values) == np.array(other.values)
        except:
            return np.array(self.values) == other

    def __ne__(self, other):
        try:
            return np.array(self.values) != np.array(other.values)
        except:
            return np.array(self.values) != other

    def __getitem__(self, arg):
        return np.array(self.values)[arg]


class SMA(Indicator):
    """ Simple moving average. """

    name = 'Simple moving average'

    def __init__(self, window):
        super(SMA, self).__init__(label='SMA ' + str(window))
        self.window = int(window)
        self.plot_type = 'primary'

    def calculate(self, candles):
        closes = np.array(candles).transpose()[2]
        values = []
        for i, price in enumerate(closes):
            if i < self.window:
                pt = sum(closes[:i+1])/(i+1)
            else:
                pt = sum(closes[i-self.window+1:i+1])/self.window
            values.append(pt)
        self.values = np.array(values)
        return self.values







class EMA(Indicator):
    """ Exponential moving average. """
    name = 'Exponential moving average'

    def __init__(self, window):
        Indicator.__init__(self, label='EMA ' + str(window))
        self.window = int(window)
        self.plot_type = 'primary'
        self.m= Decimal(2 / (self.window + 1))
        self.values = []
        
    def calculate(self, closes):
        
       
        for i, price in enumerate(closes):
            if i < self.window:
                pt = sum(closes[:i+1])/(i+1)
            else:
                pt = (closes[i] - self.values[-1]) * self.m + self.values[-1]
            self.values.append(pt)
        
        return self.values
    
    def feed(self,close):
        if len(self.values)==0:
            pt=close
        elif len(self.values)<self.window:
            pt=(self.values[-1]*len(self.values)+close)/(len(self.values)+1)
        else:
            pt = (close - self.values[-1]) * self.m + self.values[-1]
        self.values.append(pt)
    
    def getValues(self):
        return self.values
        
            

class MACD(Indicator):
    """ Moving average convergence-divergence. """

    name = 'Moving average convergence divergence'

    def __init__(self, window1, window2):
        Indicator.__init__(self, label='MACD')
        self.ma1=EMA(window1)
        self.ma2=EMA(window2)
        self.averages = sorted((self.ma1, self.ma2), key=lambda x: x.window)
        self.values=[]

    def calculate(self, closes):
        self.values = (self.averages[0].calculate(closes) -
                       self.averages[1].calculate(closes))
        return self.values
    
    def feed(self, close):
        self.ma1.feed(close)
        self.ma2.feed(close)
        self.values.append(self.averages[0].getValues()[-1]-self.averages[1].getValues()[-1])
        return self.isTrade()
    
    def isTrade(self):
        if len(self.values)>self.averages[1].window:
            if self.values[-2]>=0 and self.values[-1]<0:
                print("selling")
                return (True,"Sell")
            elif self.values[-2]<=0 and self.values[-1]>0:
                print("buying")
                return (True,"Buy")
        
        return (False,"DoNothing")
    
    def getValues(self):
        return self.values
    







