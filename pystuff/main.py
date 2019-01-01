import time

class Foo:
    done = False

    def hello(self):
        for x in range(10):
            print(str(x) + " foobar")
            import pdb; pdb.set_trace()

            time.sleep(1)
        if not self.done:
            self.done = True
            return "Hello World!"
        else:
            return "Goodbye world!"


if __name__ == "__main__":
    print(Foo().hello())
