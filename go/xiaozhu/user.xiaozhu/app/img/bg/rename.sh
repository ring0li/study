cd getup
rename 's/\.png/\.bak/' *
i=1; for x in *; do mv $x $i.png; let i=i+1; done
cd ..

cd sleep
rename 's/\.png/\.bak/' *
i=1; for x in *; do mv $x $i.png; let i=i+1; done
