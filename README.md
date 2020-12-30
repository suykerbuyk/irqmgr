A simple REST API driven server and client for monitoring and managing IRQs in a system where irqbalance may not be yielding the desired results or more fine grained control is required.

The Challenge:
Modern servers typically have many cores. Modern high performance network cards present many queues to the host, each with their own interrupt. The linux program irqbalance doesn't always evenly balance all these queue IRQs between CPUs on a multicore system. In order to make sure that a single core isn't completely overloaded by interrupt requests, it is sometimes necessary to manually set affinity to balance the IRQs. Here is some background:

    https://en.m.wikipedia.org/wiki/Network_interface_controller#Performance_and_advanced_functionality

    https://access.redhat.com/documentation/en-US/Red_Hat_Enterprise_Linux/6/html/Performance_Tuning_Guide/s-cpu-irq.html

Task: IRQ Balancing Algorithm

Below is a contrived /proc/interrupts that shows the rate of interrupts that are occuring in interrupts per day (instead of the raw count that /proc/interrupts actually shows).

Write a short program to do an approximation of the best way to evenly balance these IRQs between two CPUs.

Your script should output:

    A list of IRQs to have their affinity set to CPU0 or CPU1.

    A metric showing how closely balanced the IRQs are.

You can assume:

    Although the mock /proc/interrupts shows IRQs being serviced by both CPUs, will be pinning one interrupt to one CPU for simplicity.

The /proc/interrupt file:

       CPU0                 CPU1

132:    7805535           2676698559   IR-PCI-MSI-edge      eth0-TxRx-0
133:  177894710             78268272   IR-PCI-MSI-edge      eth0-TxRx-1
134: 3150750313             16107924   IR-PCI-MSI-edge      eth0-TxRx-2
135:  125658869             99955593   IR-PCI-MSI-edge      eth0-TxRx-3
136: 3320311515           1430447281   IR-PCI-MSI-edge      eth0-TxRx-4
137:   33721258            100610747   IR-PCI-MSI-edge      eth0-TxRx-5
138: 2707861846           1580501564   IR-PCI-MSI-edge      eth0-TxRx-6
139:   34909680             88149765   IR-PCI-MSI-edge      eth0-TxRx-7
140: 1239035616           1484418966   IR-PCI-MSI-edge      eth0-TxRx-8
141:   51448179            118527214   IR-PCI-MSI-edge      eth0-TxRx-9
142:   38185971           1941013980   IR-PCI-MSI-edge      eth0-TxRx-10
143:  132472140             72502939   IR-PCI-MSI-edge      eth0-TxRx-11
144: 3013170068           1328432100   IR-PCI-MSI-edge      eth0-TxRx-12
145:   66348784            136628241   IR-PCI-MSI-edge      eth0-TxRx-13
146: 2944162504           1412076854   IR-PCI-MSI-edge      eth0-TxRx-14
147:   32024336            108557842   IR-PCI-MSI-edge      eth0-TxRx-15
148: 1756364855           1374481202   IR-PCI-MSI-edge      eth0-TxRx-16
149:    1661862              2913153   IR-PCI-MSI-edge      eth0-TxRx-17
150: 3431403731            658237917   IR-PCI-MSI-edge      eth0-TxRx-18
151:    3071298              9025526   IR-PCI-MSI-edge      eth0-TxRx-19
152: 3445872980           1247634023   IR-PCI-MSI-edge      eth0-TxRx-20
153:    1231242              3038219   IR-PCI-MSI-edge      eth0-TxRx-21
154:  973855391           1243772811   IR-PCI-MSI-edge      eth0-TxRx-22
155:    1033227              2939261   IR-PCI-MSI-edge      eth0-TxRx-23
156: 2820232388           1187654439   IR-PCI-MSI-edge      eth0-TxRx-24
158: 1081720733           1262898017   IR-PCI-MSI-edge      eth0-TxRx-26
159:    1238642              3794156   IR-PCI-MSI-edge      eth0-TxRx-27
160:  854024413           1289247462   IR-PCI-MSI-edge      eth0-TxRx-28
161:    1185992              3543676   IR-PCI-MSI-edge      eth0-TxRx-29
162: 1034036073           1478888774   IR-PCI-MSI-edge      eth0-TxRx-30
163:    1502733              4743028   IR-PCI-MSI-edge      eth0-TxRx-31
164:       9232                25456   IR-PCI-MSI-edge      eth0
165:     189741               411172   IR-PCI-MSI-edge      eth1-TxRx-0
166:     152345               468245   IR-PCI-MSI-edge      eth1-TxRx-1
167:     157516               549207   IR-PCI-MSI-edge      eth1-TxRx-2
168:     184096               475831   IR-PCI-MSI-edge      eth1-TxRx-3
169:     153135               485629   IR-PCI-MSI-edge      eth1-TxRx-4
170:     176824               468308   IR-PCI-MSI-edge      eth1-TxRx-5
171:     146520               415216   IR-PCI-MSI-edge      eth1-TxRx-6
172:     144750               545998   IR-PCI-MSI-edge      eth1-TxRx-7
173:     147808               473153   IR-PCI-MSI-edge      eth1-TxRx-8
174:     144077               450370   IR-PCI-MSI-edge      eth1-TxRx-9
175:     146128               430602   IR-PCI-MSI-edge      eth1-TxRx-10
176:     147827               502182   IR-PCI-MSI-edge      eth1-TxRx-11
177:     144080               510764   IR-PCI-MSI-edge      eth1-TxRx-12
178:     141517               548450   IR-PCI-MSI-edge      eth1-TxRx-13
179:     139864               463763   IR-PCI-MSI-edge      eth1-TxRx-14
180:     156056               502055   IR-PCI-MSI-edge      eth1-TxRx-15
181:     148448               554651   IR-PCI-MSI-edge      eth1-TxRx-16
182:     152685               460494   IR-PCI-MSI-edge      eth1-TxRx-17
183:     137910               426358   IR-PCI-MSI-edge      eth1-TxRx-18
184:     146274               514366   IR-PCI-MSI-edge      eth1-TxRx-19
185:     152762               433858   IR-PCI-MSI-edge      eth1-TxRx-20
186:     151645               508297   IR-PCI-MSI-edge      eth1-TxRx-21
187:     159089               423519   IR-PCI-MSI-edge      eth1-TxRx-22
188:     141547               478127   IR-PCI-MSI-edge      eth1-TxRx-23
189:     150970               445886   IR-PCI-MSI-edge      eth1-TxRx-24
190:     159699               466935   IR-PCI-MSI-edge      eth1-TxRx-25
191:     142656               558960   IR-PCI-MSI-edge      eth1-TxRx-26
192:     149152               473756   IR-PCI-MSI-edge      eth1-TxRx-27
193:     149436               514896   IR-PCI-MSI-edge      eth1-TxRx-28
194:     149677               503363   IR-PCI-MSI-edge      eth1-TxRx-29
195:     144588               417862   IR-PCI-MSI-edge      eth1-TxRx-30
196:     152210               514647   IR-PCI-MSI-edge      eth1-TxRx-31

Task: IRQ Client/Server
IRQ Information Service

We would like a service that presents a restful API which allows us to:

    Get an overview of how interrupts are distributed among the different CPUs in our system. This overview should be over a specified time window, showing how many interrupts fired within the window and on which CPUs.

    Set the CPU affinity for each interrupt. This can be simplified by allowing only one CPU to be associated with any given interrupt (instead of supporting the entire mask).

    Provide a basic init script to start this service at machine boot time. Our target OS is Ubuntu 14.04 or 16.04.

    Bonus: Package the service as a basic wheel using setuptools.

You can assume:

    This service will be run as root.

    No authentication or encryption is necessary.

IRQ Service Client

We would like a client to interact with this service. It should be able to:

    Give a complete overview of the current distribution of interrupts.

    Give a summary of how many interrupts have been serviced by each CPU.

    Be able to use the service to set the affinity of various interrupts.

You can assume:

    Client can be a standalone script, packaging isn't necessary.

