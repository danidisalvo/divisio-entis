import {Component, OnInit, ViewChild} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Color} from '@angular-material-components/color-picker';
import {MatDialog} from "@angular/material/dialog";
import { ToastContainerDirective, ToastrService } from 'ngx-toastr';
import {ConfirmDialogComponent} from "../confirm-dialog/confirm-dialog.component";
import {EditNodeDialogComponent} from "../edit-node-dialog/edit-node-dialog.component";
import {environment} from '../../environments/environment';

import * as d3 from 'd3';

export type Properties = {
  [key: string]: string;
};

export interface Node {
  id: string,
  name: string;
  color: Color,
  properties: Properties,
  type: string
  children?: Array<Node>
}

@Component({
  selector: 'app-graph',
  templateUrl: './graph.component.html',
  styleUrls: ['./graph.component.css']
})

export class GraphComponent implements OnInit {
  @ViewChild(ToastContainerDirective, { static: true })
  toastContainer: ToastContainerDirective | undefined;

  private dialogWidth = '500px';
  private nodeRadius = 10
  private pathStrokeWidth = 2;
  private strokeColor = '#98989C';
  private transitionDuration = 750;

  hzoom: number
  vzoom: number

  constructor(private http: HttpClient, public dialog: MatDialog, private toastrService: ToastrService) {
    this.hzoom = 50
    this.vzoom = 50
  }

  ngOnInit(): void {
    this.toastrService.overlayContainer = this.toastContainer;

    const url = `${environment.apiUrl}`;
    this.http.get<Node>(url + "graph").subscribe({
      next: data => {
        this.drawTree(data);
      },
      error: error => {
        this.toastrService.error(error.status + ' ' + error.statusTex)
      }
    });
  }

  public clearGraph() {
    this.dialog.open(ConfirmDialogComponent, {
      autoFocus: true,
      disableClose: true,
      width: this.dialogWidth,
      data: {
        message: 'Do you want to clear the graph?',
      }
    }).afterClosed().subscribe((result: any) => {
      if (result != undefined && result.action == 'yes') {
        const url = `${environment.apiUrl}` + 'graph';
        this.http.delete(url).subscribe({
          next: () => {
            this.ngOnInit();
          },
          error: error => {
            this.toastrService.error(error.status + ' ' + error.statusTex)
          }
        });
      }
    });
  }

  public zoom() {
    this.ngOnInit();
  }

  public drawTree(treeData: any) {
    // create the nested data structure representing the tree
    const treeDataStructure: any = d3.hierarchy(treeData, (d: any) => {
      return d.children;
    });
    treeDataStructure.x0 = window.screen.height / 2;
    treeDataStructure.y0 = 0;

    const hZoomFactor = this.hzoom / 50;
    const vZoomFactor = this.vzoom / 50;
    // set the svg attributes
    d3.select("div#container").select("svg").remove();

    const svg = d3.select("div#container")
      .append("svg")
      .attr("preserveAspectRatio", "xMinYMin meet")
      .attr("viewBox", "-50 -50 " + 1800 * hZoomFactor + " " + 1800 * 2 * vZoomFactor)
      // .attr("viewBox", "-50 -50 " + window.screen.width * 2 * hZoomFactor + " " + window.screen.height * 2 * vZoomFactor)
      .classed("svg-content", true)
      .append('g');

    // create and configure the tree layout
    const treeLayout = d3.tree().size([1800, 1800]);

    // see https://stackoverflow.com/questions/63157253/how-to-initialize-a-angular-material-components-color-picker-side-ts-with-stri
    const hexToRgb = (hex: string) => {
      const shorthandRegex = /^#?([a-f\d])([a-f\d])([a-f\d])$/i;
      hex = hex.replace(shorthandRegex, (m, r, g, b) => {
        return r + r + g + g + b + b;
      });
      const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
      return result ? {
        r: parseInt(result[1], 16),
        g: parseInt(result[2], 16),
        b: parseInt(result[3], 16)
      } : null;
    }

    //////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Edit a node and add a child
    //////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    const editNode = (event: Event, d: any) => {
      // prevents the default context menu
      event.preventDefault();

      const rbg = hexToRgb(d.data.color);

      const node: Node = {
        id: d.data.id,
        name: d.data.name,
        color: rbg == null ? new Color(0, 0, 0) : new Color(rbg.r, rbg.g, rbg.b),
        type: d.data.type,
        properties: d.data.properties
      }

      const child: Node = {
        id: '',
        name: '',
        color: rbg == null ? new Color(0, 0, 0) : new Color(rbg.r, rbg.g, rbg.b),
        type: '',
        properties: {}
      }

      this.dialog.open(EditNodeDialogComponent, {
        autoFocus: true,
        disableClose: true,
        width: this.dialogWidth,
        data: {
          d: d,
          node: node,
          child: child,
          addChild: false
        }
      }).afterClosed().subscribe((result: any) => {
        if (result != undefined && result.action == 'save') {
          let url = `${environment.apiUrl}` + 'nodes';
          let body: Node;
          if (result.d.parent == null) { // if parent is null, we can only create a child
            body = {
              id: result.child.id,
              name: result.child.name,
              color: result.child.color.toHexString(),
              type: result.child.type,
              properties: result.child.properties,
              children: []
            }
          } else { // if parent is *not* null we can: a) amend the parent itself; and/or create a child
            url = `${environment.apiUrl}` + 'nodes/' + result.d.parent.data.id;
            body = { // the parent is *not* the root node
              id: result.node.id,
              name: result.node.name,
              color: result.node.color.toHexString(),
              type: result.node.type,
              properties: result.child.properties,
              children: [] // we leave them empty because the backend will not set them
            }
            if (result.addChild) {
              body.children?.push({
                id: result.child.id,
                name: result.child.name,
                color: result.child.color.toHexString(),
                type: result.child.type,
                properties: result.child.properties,
                children: []
              })
            }
          }
          const httpOptions = {
            headers: new HttpHeaders({
              'Content-Type': 'application/json'
            })
          };
          this.http.put(url, body, httpOptions).subscribe({
            next: () => {
              // amend the current node
              result.d.data.name = result.node.name
              result.d.data.color = result.node.color.toHexString();
              result.d.data.type = result.node.type

              // add a child
              if (result.addChild) {
                // show all hidden children, so we can see the new node after it has been created
                if (result.d._children != null) {
                  result.d.children = result.d._children;
                  result.d._children = null;
                }
                // create the children array if it does not exist
                if (result.d.children == null) {
                  result.d.children = [];
                }
                const child = {
                  depth: result.d.depth + 1,
                  height: 0,
                  children: null,
                  _children: null,
                  parent: result.d,
                  data: {
                    id: result.child.id,
                    name: result.child.name,
                    color: result.child.color.toHexString(),
                    type: result.child.type,
                    properties: result.child.properties,
                    children: []
                  }
                };
                result.d.children.push(child);
              }

              if (result.parent !== result.d.parent.data.id) {
                let url = `${environment.apiUrl}` + 'nodes/' + result.d.parent.data.id + '/' + result.d.data.id + '/' + result.targetNode;
                this.http.post(url, body, httpOptions).subscribe({
                  next: data => {
                    this.drawTree(data);
                  },
                  error: error => {
                    this.toastrService.error(error.status + ' ' + error.statusTex)
                  }
                })
              } else {
                update(result.d);
              }
            },
            error: error => {
              this.toastrService.error(error.status + ' ' + error.statusTex)
            }
          });
        }
      });

      // required when replacing the default context menu
      return false;
    }

    //////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Delete a node
    //////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    const deleteNode = (event: Event, d: any) => {
      this.dialog.open(ConfirmDialogComponent, {
        autoFocus: true,
        disableClose: true,
        width: this.dialogWidth,
        data: {
          message: 'Do you want to delete "' + d.data.name + '" and its children?',
          d: d
        }
      }).afterClosed().subscribe((result: any) => {
        if (result != undefined && result.action == 'yes') {
          const url = `${environment.apiUrl}` + 'nodes/' + result.d.parent.data.id + '/' + result.d.data.id;
          this.http.delete(url).subscribe({
            next: () => {
              const children: any = [];
              d.parent.children.forEach((child: any) => {
                if (child.id != d.id) {
                  children.push(child);
                }
              });
              d.parent.children = children.length == 0 ? null : children;
              update(d.parent)
            },
            error: error => {
              this.toastrService.error(error.status + ' ' + error.statusTex)
            }
          });
        }
      });
    }

    //////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Updates tree's nodes and links
    //////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    const update = (source: any) => {

      const hierarchyNode = treeLayout(treeDataStructure);
      const descendants = hierarchyNode.descendants();

      descendants.forEach(function (d) {
        d.y = d.depth * 180
      });

      let i = 0;

      // start updating the nodes
      const nodes = svg.selectAll('g.node')
        .data(descendants, (d: any) => {
          return d.id || (d.id = ++i);
        });

      // enter the new nodes at their parent's previous positions
      const nodeEnter = nodes.enter().append('g')
        .attr('class', 'node')
        .attr('transform', () => {
          return 'translate(' + source.y0 + ',' + source.x0 + ')';
        })
        .on('click', toggleChildren)
        .on('contextmenu', editNode);

      // add the circle
      nodeEnter.append('circle')
        .attr('class', 'node')
        .attr('r', 0)
        .style('fill', (d: any) => {
          return d.data.color;
        })
        .attr('cursor', (d: any) => {
          return d._children ? 'pointer' : 'auto';
        });

      // add the label
      nodeEnter.append('text')
        .attr('dy', '.35em')
        .attr('x', (d: any) => {
          return d.children || d._children ? -13 : 13;
        })
        .attr('text-anchor', (d: any) => {
          return d.children || d._children ? 'end' : 'start';
        })
        .text((d: any) => {
          return d.data.name.startsWith(' ') ?  '' : d.data.name;
        });

      // @ts-ignore
      const nodeUpdate = nodeEnter.merge(nodes);

      // transition to the proper position for the node
      nodeUpdate.transition()
        .duration(this.transitionDuration)
        .attr('transform', (d: any) => {
          return 'translate(' + d.y + ',' + d.x + ')';
        });

      // update attributes and style
      nodeUpdate.select('circle.node')
        .attr('r', this.nodeRadius)
        .style('fill', (d: any) => {
          return d.data.color;
        })
        .attr('cursor', (d: any) => {
          return d._children ? 'pointer' : 'auto';
        });

      nodeUpdate.select('text')
        .attr('dy', '.35em')
        .attr('x', (d: any) => {
          return d.children || d._children ? -13 : 13;
        })
        .attr('text-anchor', (d: any) => {
          return d.children || d._children ? 'end' : 'start';
        })
        .text((d: any) => {
          return d.data.name;
        });

      // remove the exiting nodes
      const nodeExit = nodes.exit()
        .transition()
        .duration(this.transitionDuration)
        .attr('transform', () => {
          return 'translate(' + source.y + ',' + source.x + ')';
        })
        .remove();

      // set the circle radius to 0
      nodeExit.select('circle').attr('r', 0);

      // set text opacity to 0
      nodeExit.select('text').style('fill-opacity', 0);

      // start updating the links
      const links = hierarchyNode.descendants().slice(1);
      const link = svg.selectAll('path.link')
        .data(links, (d: any) => {
          return d.id;
        });

      // enter new links at their parents' previous positions
      const linkEnter = link.enter().insert('path', 'g')
        .attr('class', 'link')
        .attr('fill', 'none')
        .attr('stroke', this.strokeColor)
        .attr('stroke-width', this.pathStrokeWidth)
        .attr('d', () => {
          const src = {x: source.x0, y: source.y0};
          const dst = {x: source.x0, y: source.y0};
          return diagonal(src, dst)
        })
        .on('click', deleteNode);

      // @ts-ignore
      const linkUpdate = linkEnter.merge(link);

      // transition back to the parent element position
      linkUpdate.transition()
        .duration(this.transitionDuration)
        .attr('d', (d: any) => {
          return diagonal(d, d.parent)
        });

      // remove the exiting links
      link.exit().transition()
        .duration(this.transitionDuration)
        .attr('d', () => {
          const src = {x: source.x, y: source.y}
          const dst = {x: source.x, y: source.y}
          return diagonal(src, dst)
        })
        .remove();

      // store the old positions for transition
      descendants.forEach((d: any) => {
        d.x0 = d.x;
        d.y0 = d.y;
      });

      // create the path
      function diagonal(s: any, d: any) {
        return `M ${s.y} ${s.x}
            C ${(s.y + d.y) / 2} ${s.x},
              ${(s.y + d.y) / 2} ${d.x},
              ${d.y} ${d.x}`
      }

      // toggle children visibility and trigger an update
      function toggleChildren(event: any, d: any) {
        if (d.children) {
          d._children = d.children;
          d.children = null;
        } else {
          d.children = d._children;
          d._children = null;
        }
        update(d);
      }
    }

    update(treeDataStructure);
  }
}
