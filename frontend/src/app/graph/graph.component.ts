import {Component, OnInit} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Color} from '@angular-material-components/color-picker';
import {MatDialog} from "@angular/material/dialog";
import {ConfirmDialogComponent} from "../confirm-dialog/confirm-dialog.component";
import {EditNodeDialogComponent} from "../edit-node-dialog/edit-node-dialog.component";
import {environment} from '../../environments/environment';

import * as d3 from 'd3';
import { v4 as uuidv4 } from 'uuid';

interface Node {
  name: string;
  readOnly: boolean,
  color: Color,
  division: string,
  notes: string
  children?: Array<Node>
}

@Component({
  selector: 'app-graph',
  templateUrl: './graph.component.html',
  styleUrls: ['./graph.component.css']
})
export class GraphComponent implements OnInit {

  private dialogWidth = '500px';
  private nodeRadius = 10
  private pathStrokeWidth = 2;
  private strokeColor = '#98989C';
  private margin = {top: 20, right: 50, bottom: 20, left: 50};
  private height = 1080 - this.margin.top - this.margin.bottom;
  private width = 1920 - this.margin.left - this.margin.right;
  private transitionDuration = 750;

  message!: string | null;

  constructor(private http: HttpClient, public dialog: MatDialog) {
  }

  ngOnInit(): void {
    const url = `${environment.apiUrl}`;
    this.http.get<Node>(url + "graph").subscribe({
      next: data => {
        this.message = null;
        this.drawTree(data);
      },
      error: error => {
        this.message = error.error.status + ': ' + error.error.message;
      }
    });
  }

  private drawTree(treeData: any) {
    // create the nested data structure representing the tree
    const treeDataStructure: any = d3.hierarchy(treeData, (d: any) => {
      return d.children;
    });
    treeDataStructure.x0 = this.height / 2;
    treeDataStructure.y0 = 0;

    // set the svg attributes
    const svg = d3.select('svg')
      .attr('width', this.width + this.margin.right + this.margin.left)
      .attr('height', this.height + this.margin.top + this.margin.bottom)
      .append('g')
      .attr('transform', 'translate(' + this.margin.left + ',' + this.margin.top + ')');

    // create and configure the tree layout
    const treeLayout = d3.tree().size([this.height, this.width]);

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
        name: d.data.name,
        readOnly: d.data.readOnly,
        color: rbg == null ? new Color(0, 0, 0) : new Color(rbg.r, rbg.g, rbg.b),
        division: d.data.division,
        notes: d.data.notes
      }

      const child: Node = {
        name: '',
        readOnly: false,
        color: new Color(0, 0, 0),
        division: '',
        notes: ''
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
              name: result.child.name,
              readOnly: false,
              color: result.child.color.toHexString(),
              division: result.child.division,
              notes: result.child.notes,
              children: []
            }
          } else { // if parent is *not* null we can: a) amend the parent itself; and/or create a child
            url = `${environment.apiUrl}` + 'nodes/' + result.d.parent.data.name;
            body = { // the parent is *not* the root node
              name: result.node.name,
              readOnly: result.node.readOnly,
              color: result.node.color.toHexString(),
              division: result.node.division,
              notes: result.node.notes,
              children: []
            }
            if (result.addChild) {
              body.children?.push({
                name: result.child.name,
                readOnly: false,
                color: result.child.color.toHexString(),
                division: result.child.division,
                notes: result.child.notes,
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
              this.message = null;

              // amend the current node
              result.d.data.color = result.node.color.toHexString();
              result.d.data.notes = result.node.notes;
              result.d.data.division = result.node.division;

              // add a child
              if (result.addChild) {
                // show all hidden children, so we can see the new node after it has been created
                if (result.d._children != null) {
                  result.d.children = result.d._children;
                  result.d._children = null;
                }
                // create the children array if does not exist
                if (result.d.children == null) {
                  result.d.children = [];
                }
                const child = {
                  id: uuidv4(),
                  depth: result.d.depth + 1,
                  height: 0,
                  children: null,
                  _children: null,
                  parent: result.d,
                  data: {
                    name: result.child.name,
                    color: result.child.color.toHexString(),
                    notes: result.child.notes,
                    division: result.child.division,
                    readOnly: false,
                    children: []
                  }
                };
                result.d.children.push(child);
              }
              update(result.d);
            },
            error: error => {
              console.error(error);
              this.message = error.error.status + ': ' + error.error.message;
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
          const url = `${environment.apiUrl}` + 'nodes/' + result.d.parent.data.name + '/' + result.d.data.name;
          this.http.delete(url).subscribe({
            next: () => {
              this.message = null;
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
              console.error(error);
              this.message = error.error.status + ': ' + error.error.message;
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
      const nodes = hierarchyNode.descendants();

      nodes.forEach(function (d) {
        d.y = d.depth * 180
      });

      let i = 0;

      // start updating the nodes
      const node = svg.selectAll('g.node')
        .data(nodes, (d: any) => {
          return d.id || (d.id = ++i);
        });

      // enter the new nodes at their parent's previous positions
      const nodeEnter = node.enter().append('g')
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
          return d.data.name;
        });

      // @ts-ignore
      const nodeUpdate = nodeEnter.merge(node);

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

      // remove the exiting nodes
      const nodeExit = node.exit()
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
      nodes.forEach((d: any) => {
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
