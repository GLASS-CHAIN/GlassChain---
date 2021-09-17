# paracross   

## 

 ， ， 。

 
 1. : 
 1.  :  
 1. paracros ： 
 1. paracros ： 

 
 1. A: conis/token -> paracros 
 1. A: paracros  -> paracros 
 1. A: paracros  -> 
 1. B A): paracros  -> paracros 

## 

 ， 
 1. ，  title  user.p.guodun.paracross， 

 
 1. ， 
 1. bitmap，   bit 

## 

asset-transfer ， ， 


  transfer
 * 
   1. paracros ， balance -
   1. paracros ， balance +
 *  ， )
   1.  paracros   balance +

  withdraw
 * 
   1.  paracros   balance -
 *  ， )
   1. commit 
   1. paracros ， balance -
   1. paracros ， balance +

 <-  cross-transfer
>cross-transfer transfe withdra transfer transfe withdraw

 　=　assetExec + assetSymbol 
  1. ：coins+BTY,token+CCNY
  1. :user.p.test.coins + FZM,
  1. paracros : ：paracross　+ user.p.test.coins.FZM， : user.p.test.paracross + coins.BTY
  1.   
  1. titl transfe 
 :
```
				exec                    symbol                              tx.title=user.p.test1   tx.title=user.p.test2
1. ：
				coins                   bty                                 transfer                 transfer
				paracross               user.p.test1.coins.fzm              withdraw                 transfer

2. ：
				user.p.test1.coins      fzm                                 transfer                 NAN
                user.p.test1.paracross  coins.bty                           withdraw                 NAN
                user.p.test1.paracross  paracross.user.p.test2.coins.cny    withdraw                 NAN

 user.p.test1.paracross.paracross.user.p.test2.coins.cn ：
user.p.test1.paracross paracros ，　paracross.user.p.test2.coins.cn paracros paracros user.p.test2.coins.cn 
```

  
 1. 
 1. 
 1. 

### kv ， 

 1. kv 
    1. 
    1.    ， )
    1. 
 1. 
    1. 
 1. 
    1. ：  

```
                                                            
account                              A    B
1  5bty    5           0                5              0               0       0           0
 ar 
2  5bty    10          0                5              5               0       0           0
 ar 
3  4bty    10          4                1              5               0       0           0      
         10          4                1              5               4       4           0      
4  3bty      10          4                1              5               1       1           0
 
5   2bty   10          4                1              5               3       1           2
 ， 2bt 
6       10          4                1              5               3       1           2       
  1bty      10          4                1              5               2       1           1       
              10           3                1              6               2       1           1       
```

### <-  cross-transfer 
```
# Alice ５coins-bty -> user.p.test. :

                    coins       paracross:Addr(Alice)   paracross:Addr(user.p.test.paracross)    user.p.test.paracross-coins-bty:Addr(Alice) 
1 Alice                5
2 t                 0　　　　　　　　 5       
3 cross-transfer       0            5-5=0                   0+5=5                                          0+5=5

# Alice ５paracross-coins.bty -> 
4 cross-transfer                    　5                   5-5=0                                       5-5=0
5 withdraw           　5              0

# Bob 5 user.p.test.coins.fzm -> 
                    paracross-user.p.test.coins.fzm:Addr(Bob)    user.p.test.coins.fzm      user.p.test.paracross:Addr(Bob)   user.p.test.paracross:Addr(paracross)
1 Bob                                                                       5
2 to paracros 　　            　　　　　　　　                               0                       5       
3 cross-transfer                  0+5=5                                                             5-5=0                             0+5=5     

# Bob ５exec:paracross　symbol:user.p.test.coins.fzm -> 
4 cross-transfer                  5-5=0                                                             0+5=5                                5-5=0
5 withdraw                                                                  5                       5-5=0


```