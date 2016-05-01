package goperiodicity

func scale(old_min, old_max, new_min, new_max, value float64) float64{
    return ((value - old_min) / (old_max - old_min) ) * (new_max - new_min) + new_min
}
